package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sugam/golang-react-todo/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var collection *mongo.Collection

func init() {
	loadTheEnv()
	CreateDBInstance()
}
func CreateDBInstance() {
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	clientOptions := options.Client().ApplyURI(connectionString)
	fmt.Println(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(clientOptions)
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mogodb...")
	client.Database(dbName).Collection(collectionName)
	fmt.Println("collection instance created...")
}

func loadTheEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading in .env file :: %v", err)
	}
}
func (server *Server) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//payload := getAllTasks()
	payload, err := models.FindAll(server.DB)
	fmt.Println("payload", payload)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(payload)
}

func (server *Server) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Contol-Allow-Headers", "Content-Type")
	var task models.ToDoList
	//task.ID = bson.NewObjectId()
	json.NewDecoder(r.Body).Decode(&task)
	fmt.Println("task", task)
	payload, err := models.CreateTodoList(server.DB, task)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("start....67")
	json.NewEncoder(w).Encode(payload)
}

func (server *Server) TaskComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Contol-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	taskComplete(params["_id"])
	json.NewEncoder(w).Encode(params["_id"])
}

func (server *Server) UndoTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Contol-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	undoTask(params["_id"])
	json.NewEncoder(w).Encode(params["_id"])
}

func (server *Server) DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Contol-Allow-Headers", "Content-Type")
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		return
	}
	fmt.Println(pid)
	_, err = models.DeleteTodoList(server.DB, uint(pid))
	if err != nil {
		json.NewEncoder(w).Encode("unable to delete")
	}
}

func (server *Server) DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Contol-Allow-Headers", "Content-Type")
	count := deleteAllTask()
	json.NewEncoder(w).Encode(count)
}

// func getAllTasks() []primitive.M {
// 	fmt.Println("-------------------")
// 	cur, err := collection.Find(context.Background(), bson.D{{}})
// 	fmt.Println(cur)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(cur)
// 	fmt.Println("-------------------1")
// 	var results []primitive.M
// 	for cur.Next(context.Background()) {
// 		var result bson.M
// 		err = cur.Decode(&result)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Println("-------------------2")
// 		results = append(results, primitive.M(result))
// 	}
// 	if err := cur.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("-------------------2")
// 	cur.Close(context.Background())
// 	return results
// }

func taskComplete(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count::", result.ModifiedCount)
}

func insertOneTask(task models.ToDoList) {
	fmt.Println("db::", task)
	result, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		fmt.Println("Getting error")
		log.Fatal("error occured", err)
	}
	fmt.Println("inserted data succed on this  id::", result.InsertedID)
}

func undoTask(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"status": false}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("updated to undo list :: ", result.ModifiedCount)
}

// func deleteTask(task string) {
// 	id, _ := primitive.ObjectIDFromHex(task)
// 	filter := bson.M{"id": id}
// 	result, err := collection.DeleteOne(context.Background(), filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Deleted succssfully ::", result.DeletedCount)
// }

func deleteAllTask() int {
	result, err := collection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	return int(result.DeletedCount)
}
