package main

import (
	"log"
	"net/http"
	"os"
	"storeSystem/internal/database"
	"storeSystem/internal/handlers"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://storageuser:storagepass@localhost:5433/storagedb?sslmode=disable"
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	log.Printf("starting server on port %s", serverPort)
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("успешно подключено к бд")

	itemStore := database.NewItemStore(db)
	handler := handlers.NewHandlers(itemStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/items", methodHandler(handler.GetAllItems, "GET"))
	mux.HandleFunc("/items/create", methodHandler(handler.CreateItem, "POST"))
	mux.HandleFunc("/items/", itemIDHandler(handler))

	loggedMux := loggingMiddleware(mux)
	//cors middleware

	serverAddr := ":" + serverPort

	err = http.ListenAndServe(serverAddr, loggedMux)

	if err != nil {
		log.Fatal(err)
	}
}

func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}

func itemIDHandler(handler *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetItemByID(w, r)
		case http.MethodPut:
			handler.UpdateItem(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s%s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
