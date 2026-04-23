package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type ReportConfig struct {
	ReportType string `json:"report_type"`
	DateFrom   string `json:"date_from"`
	DateTo     string `json:"date_to"`
	//clientID   int
	//date       string
}

func (h *Handlers) GenerateStockReport(config *ReportConfig) error {
	stocks, err := h.stockStore.GetAll()
	if err != nil {
		return fmt.Errorf("Ошибка получения остатков")
	}

	file, err := os.Create("stock.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	//header := []string{"ID", "Name", "Age"}
	//if err := writer.Write(header); err != nil {
	//	log.Fatal(err)
	//}

	for _, stock := range stocks {
		jsonData, err := json.Marshal(stock)
		fmt.Println(string(jsonData))

		var result map[string]interface{}
		err = json.Unmarshal(jsonData, &result)

		if err != nil {
			fmt.Printf("%v\n", err)
			return err
		}

		values := make([]string, 0, len(result))
		for _, value := range result {
			values = append(values, fmt.Sprint(value))
		}

		fmt.Println(values)
		if err := writer.Write(values); err != nil {
			return err
		}
	}

	//filename := config.reportType
	//f, err := file.Open()
	//if err != nil {
	//
	//}
	//fileData := helpers.FileDataType{
	//	FileName: filename,  // Имя файла
	//	Data:     fileBytes, // Содержимое файла в виде байтового среза
	//}
	return err
}

//
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

func (h *Handlers) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var config *ReportConfig

	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(*config)

	var generationError error
	switch config.ReportType {
	case "stock":
		generationError = h.GenerateStockReport(config)
	default:
		generationError = fmt.Errorf("unknown report type")
	}
	if generationError != nil {
		respondWithError(w, http.StatusInternalServerError, generationError.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, "")
}
