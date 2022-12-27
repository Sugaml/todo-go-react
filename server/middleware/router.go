package middleware

import "github.com/gorilla/mux"

func (server *Server) Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", server.GetAllTasks).Methods("GET")
	router.HandleFunc("/api/task", server.CreateTask).Methods("POST")
	router.HandleFunc("/api/tasks/{id}", server.TaskComplete).Methods("PUT")
	router.HandleFunc("/api/task-undo/{id}", server.UndoTask).Methods("PUT")
	router.HandleFunc("/api/delete-task/{id}", server.DeleteTask).Methods("DELETE")
	router.HandleFunc("/api/delete-all-task", server.DeleteAllTasks).Methods("Delete")
	return router
}
