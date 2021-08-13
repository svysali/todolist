package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/svysali/todolist/db"
	"github.com/svysali/todolist/db/sqlc"
)

var (
	queries *sqlc.Queries
)

type ErrorCode int

type Error struct {
	Message    string    `json:"message"`
	StatusCode int       `json:"code"`
	ErrorCode  ErrorCode `json:"error_code,omitempty"`
}

type Todoitem struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"app":"todolist", "status":"ok"}`)
}

func listItems(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	items, err := queries.ListItems(ctx)
	if err != nil {
		log.Error("Could not fetch items: %w", err)
		RespondWithError(w, r, "Could not fetch items", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	RespondWithStatus(w, r, items, 200)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	item := sqlc.Item{}
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		RespondWithError(w, r, "bad request", 500)
		return
	}
	log.WithFields(log.Fields{"title": item.Title}).Info("Add new TodoItem. Saving to database.")
	created, err := queries.CreateItem(ctx, item.Title)
	if err != nil {
		log.Error("Could not create item: %w", err)
		RespondWithError(w, r, "Could not create item", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	RespondWithStatus(w, r, created, 200)
}

func RespondWithError(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	apiErr := &Error{
		Message:    msg,
		StatusCode: statusCode,
	}
	RespondWithStatus(w, r, apiErr, statusCode)
}

func RespondWithStatus(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)
	return
}

func main() {
	log.Info("Starting Todolist Server")
	configureDatabase()
	queries = sqlc.New(db.GetDBConn())

	router := mux.NewRouter()
	router.HandleFunc("/healthz", healthz).Methods("GET")
	router.HandleFunc("/list", listItems).Methods("GET")
	router.HandleFunc("/todo", createItem).Methods("POST")
	http.ListenAndServe(":8080", router)
}

func configureDatabase() {
	err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}
}
