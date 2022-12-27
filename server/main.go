package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/sugam/golang-react-todo/db/postgres"
	"github.com/sugam/golang-react-todo/middleware"
)

func main() {
	fmt.Println("start..")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not coming through %v", err)
	}

	DBDriver := os.Getenv("DB_DRIVER")
	fmt.Println(DBDriver)
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	fmt.Println(DBURL)
	conn, err := gorm.Open(DBDriver, DBURL)
	if err != nil {
		log.Printf("Error to connect DB...%v", err)
	}
	server := middleware.NewServer(conn)
	postgres.Migration(conn)
	r := server.Router()
	fmt.Println("starting the server on port 9000 ...")
	log.Fatal(http.ListenAndServe(":9000", r))
}
