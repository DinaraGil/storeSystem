package main

import (
	"log"
	"net/http"
	"os"
	"storeSystem/internal/database"
	"storeSystem/internal/handlers"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	workerStore := database.NewWorkerStore(db)
	deliveryListStore := database.NewDeliveryListStore(db)
	deliveryStore := database.NewDeliveryStore(db)
	handler := handlers.NewHandlers(itemStore, workerStore, deliveryListStore, deliveryStore)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/items", func(r chi.Router) {
		r.Get("/", handler.GetAllItems)     // GET /items
		r.Post("/", handler.CreateItem)     // POST /items
		r.Get("/{id}", handler.GetItemByID) // GET /items/1
		r.Put("/{id}", handler.UpdateItem)  // PUT /items/1
	})

	router.Route("/delivery_lists", func(r chi.Router) {
		r.Get("/", handler.GetAllDeliveryLists)
		r.Post("/", handler.CreateDeliveryList)
		r.Get("/{id}", handler.GetDeliveryListByID)
	})

	router.Route("/deliveries", func(r chi.Router) {
		r.Get("/", handler.GetAllDeliveries)
		r.Post("/", handler.CreateDelivery)
		r.Get("/{id}", handler.GetDeliveryByID)
		r.Put("/{id}", handler.UpdateDelivery)
	})

	router.Route("/workers", func(r chi.Router) {
		r.Post("/", handler.CreateWorker) // POST /workers
	})

	router.Post("/auth/login", handler.Login)

	//cors middleware

	serverAddr := ":" + serverPort

	err = http.ListenAndServe(serverAddr, router)
	if err != nil {
		log.Fatal(err)
	}
}
