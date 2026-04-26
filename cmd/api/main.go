package main

import (
	"log"
	"net/http"
	"os"
	"storeSystem/internal/config"
	"storeSystem/internal/database"
	"storeSystem/internal/handlers"
	"storeSystem/internal/minio"
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

	config.LoadConfig()

	minioClient := minio.NewMinioClient()
	err := minioClient.InitMinio()
	if err != nil {
		log.Fatalf("Ошибка инициализации Minio: %v", err)
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
	counterpartyStore := database.NewCounterpartyStore(db)
	stockStore := database.NewStockStore(db)
	reportStore := database.NewReportStore(db)

	handler := handlers.NewHandlers(
		itemStore,
		workerStore,
		deliveryListStore,
		deliveryStore,
		counterpartyStore,
		stockStore,
		minioClient,
		reportStore,
	)

	go handler.ListenEvents(databaseURL)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Post("/auth/login", handler.Login)

	router.Route("/workers", func(r chi.Router) {
		r.Post("/", handler.CreateWorker) // POST /workers
	})

	router.Group(func(router chi.Router) {
		router.Use(handler.AuthMiddleware)

		router.Route("/items", func(r chi.Router) {
			r.Get("/", handler.GetAllItems)                               // GET /items
			r.With(handlers.RequireAdmin()).Post("/", handler.CreateItem) // POST /items
			r.Get("/{id}", handler.GetItemByID)                           // GET /items/1
			r.Put("/{id}", handler.UpdateItem)                            // PUT /items/1
		})

		router.Route("/delivery_lists", func(r chi.Router) {
			r.Get("/", handler.GetAllDeliveryLists)
			r.With(handlers.RequireAdmin()).Post("/", handler.CreateDeliveryList)
			r.Get("/{id}", handler.GetDeliveryListByID)
			r.With(handlers.RequireAdmin()).Post("/upload", handler.UploadDeliveryList)
		})

		router.Route("/deliveries", func(r chi.Router) {
			r.Get("/", handler.GetAllDeliveries)
			r.With(handlers.RequireAdmin()).Post("/", handler.CreateDelivery)
			r.Get("/{id}", handler.GetDeliveryByID)
			r.Put("/{id}", handler.UpdateDelivery)
			r.Get("/{id}/lists", handler.GetDeliveryListsByDeliveryID)
			router.Put("/{id}/complete", handler.CompleteDelivery)
		})

		router.Route("/auth", func(r chi.Router) {
			r.Get("/me", handler.Me)
			r.Post("/logout", handler.Logout)
		})

		router.Get("/ws/deliveries/{delivery_id}/scanners/{scanner_id}", handler.ScanSocket)

		router.With(handlers.RequireAdmin()).Route("/counterparties", func(r chi.Router) {
			r.Get("/", handler.GetAllCounterparties)
			r.Get("/{id}", handler.GetCounterpartyByID)
			r.Post("/", handler.CreateCounterparty)
		})

		router.With(handlers.RequireAdmin()).Get("/stocks", handler.GetAllStocks)
		
		router.With(handlers.RequireAdmin()).Post("/reports/new", handler.GenerateReport)
		router.With(handlers.RequireAdmin()).Get("/reports", handler.GetUsersReports)
	})

	//cors middleware

	serverAddr := ":" + serverPort

	err = http.ListenAndServe(serverAddr, router)
	if err != nil {
		log.Fatal(err)
	}
}
