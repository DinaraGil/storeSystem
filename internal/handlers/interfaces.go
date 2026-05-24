package handlers

import (
	"storeSystem/internal/database"
	"storeSystem/internal/helpers"
	"storeSystem/internal/models"
)

type ItemStore interface {
	GetAll() ([]models.Item, error)
	GetByID(id int) (*models.Item, error)
	Create(input models.CreateItemInput) (*models.Item, error)
	Update(id int, input models.UpdateItemInput) (*models.Item, error)
}

type StockStore interface {
	GetAll() ([]models.Stock, error)
}

type CounterpartyStore interface {
	GetAll() ([]models.Counterparty, error)
	GetByID(id int) (*models.Counterparty, error)
	Create(input models.CreateCounterpartyInput) (*models.Counterparty, error)
}

type DeliveryStore interface {
	GetAll() ([]models.Delivery, error)
	GetErrorDeliveries() ([]models.Delivery, error)
	GetByID(id int) (*models.Delivery, error)
	Create(input models.CreateDeliveryInput) (*models.Delivery, error)
	CompleteDelivery(id int) error
}

type DeliveryListStore interface {
	GetAll() ([]models.DeliveryList, error)
	GetByID(id int) (*models.DeliveryList, error)
	GetByDeliveryID(id int) ([]models.DeliveryList, error)
	Create(input models.CreateDeliveryListInput) (*models.DeliveryList, error)
	ProcessScannerEvent(deliveryID int, evt models.Event, workerID int) (*database.DeliveryListUpdateDTO, error)
}

type ShipmentStore interface {
	GetAll() ([]models.Shipment, error)
	GetErrorShipments() ([]models.Shipment, error)
	GetByID(id int) (*models.Shipment, error)
	Create(input models.CreateShipmentInput) (*models.Shipment, error)
	CompleteShipment(shipmentID int) error
}

type ShipmentListStore interface {
	GetAll() ([]models.ShipmentList, error)
	GetByID(id int) (*models.ShipmentList, error)
	GetByShipmentID(id int) ([]models.ShipmentList, error)
	Create(input models.CreateShipmentListInput) (*models.ShipmentList, error)

	ProcessScannerEvent(
		shipmentID int,
		evt models.Event,
		workerID int,
	) (*database.ShipmentListUpdateDTO, error)
}

type WorkerStore interface {
	GetByUsername(username string) (*models.Worker, error)
	Create(input models.CreateWorkerInput) (*models.Worker, error)
}

type MinioService interface {
	CreateOne(data helpers.FileDataType) (*helpers.UploadedFile, error)
	GetOne(objectID string) (string, error)
}

type ReportStore interface {
	Create(userID int, reportType, filename, objectID, bucket, dateFrom, dateTo string) error
	GetByUserID(userID int) ([]models.Report, error)
}

type Conn interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	Close() error
}
