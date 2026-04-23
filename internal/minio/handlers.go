package minio

//
//import (
//	"io"
//	"net/http"
//	"storeSystem/internal/helpers"
//
//	"github.com/minio/minio-go/v7"
//)
//
//type Handler struct {
//	minioService minio.Client
//}
//
//func NewMinioHandler(
//	minioService minio.Client,
//) *Handler {
//	return &Handler{
//		minioService: minioService,
//	}
//}
//
//type Services struct {
//	minioService minio.Client // Сервис у нас только один - minio, мы планируем его использовать, поэтому передаем
//}
//
//// Handlers структура всех хендлеров, которые используются для обозначения действия в роутах
//type Handlers struct {
//	minioHandler Handler // Пока у нас только один роут
//}
//
//// NewHandler создает экземпляр Handler с предоставленными сервисами
//func NewHandler(
//	minioService Client,
//) (*Services, *Handlers) {
//	return &Services{
//			minioService: minioService,
//		}, &Handlers{
//			// инициируем Minio handler, который на вход получает minio service
//			minioHandler: *NewMinioHandler(minioService),
//		}
//}
//
//// CreateOne обработчик для создания одного объекта в хранилище MinIO из переданных данных.
//func (h *Handler) CreateOne(c *gin.Context) {
//	// Получаем файл из запроса
//	file, err := c.FormFile("file")
//	if err != nil {
//		// Если файл не получен, возвращаем ошибку с соответствующим статусом и сообщением
//		c.JSON(http.StatusBadRequest, errors.ErrorResponse{
//			Status:  http.StatusBadRequest,
//			Error:   "No file is received",
//			Details: err,
//		})
//		return
//	}
//
//	// Открываем файл для чтения
//	f, err := file.Open()
//	if err != nil {
//		// Если файл не удается открыть, возвращаем ошибку с соответствующим статусом и сообщением
//		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
//			Status:  http.StatusInternalServerError,
//			Error:   "Unable to open the file",
//			Details: err,
//		})
//		return
//	}
//	defer f.Close() // Закрываем файл после завершения работы с ним
//
//	// Читаем содержимое файла в байтовый срез
//	fileBytes, err := io.ReadAll(f)
//	if err != nil {
//		// Если не удается прочитать содержимое файла, возвращаем ошибку с соответствующим статусом и сообщением
//		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
//			Status:  http.StatusInternalServerError,
//			Error:   "Unable to read the file",
//			Details: err,
//		})
//		return
//	}
//
//	// Создаем структуру FileDataType для хранения данных файла
//	fileData := helpers.FileDataType{
//		FileName: file.Filename, // Имя файла
//		Data:     fileBytes,     // Содержимое файла в виде байтового среза
//	}
//
//	// Сохраняем файл в MinIO с помощью метода CreateOne
//	link, err := h.minioService.CreateOne(fileData)
//	if err != nil {
//		// Если не удается сохранить файл, возвращаем ошибку с соответствующим статусом и сообщением
//		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
//			Status:  http.StatusInternalServerError,
//			Error:   "Unable to save the file",
//			Details: err,
//		})
//		return
//	}
//
//	// Возвращаем успешный ответ с URL-адресом сохраненного файла
//	c.JSON(http.StatusOK, responses.SuccessResponse{
//		Status:  http.StatusOK,
//		Message: "File uploaded successfully",
//		Data:    link, // URL-адрес загруженного файла
//	})
//}
//
//// GetOne обработчик для получения одного объекта из бакета Minio по его идентификатору.
//func (h *Handler) GetOne(c *gin.Context) {
//	// Получаем идентификатор объекта из параметров URL
//	objectID := c.Param("objectID")
//
//	// Используем сервис MinIO для получения ссылки на объект
//	link, err := h.minioService.GetOne(objectID)
//	if err != nil {
//		// Если произошла ошибка при получении объекта, возвращаем ошибку с соответствующим статусом и сообщением
//		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
//			Status:  http.StatusInternalServerError,
//			Error:   "Enable to get the object",
//			Details: err,
//		})
//		return
//	}
//
//	// Возвращаем успешный ответ с URL-адресом полученного файла
//	c.JSON(http.StatusOK, responses.SuccessResponse{
//		Status:  http.StatusOK,
//		Message: "File received successfully",
//		Data:    link, // URL-адрес полученного файла
//	})
//}
//
//// RegisterRoutes - метод регистрации всех роутов в системе
//func (h *Handlers) RegisterRoutes(router *gin.Engine) {
//
//	// Здесь мы обозначили все эндпоинты системы с соответствующими хендлерами
//	minioRoutes := router.Group("/files")
//	{
//		minioRoutes.POST("/", h.minioHandler.CreateOne)
//		minioRoutes.POST("/many", h.minioHandler.CreateMany)
//
//		minioRoutes.GET("/:objectID", h.minioHandler.GetOne)
//		minioRoutes.GET("/many", h.minioHandler.GetMany)
//
//		minioRoutes.DELETE("/:objectID", h.minioHandler.DeleteOne)
//		minioRoutes.DELETE("/many", h.minioHandler.DeleteMany)
//	}
//
//}
