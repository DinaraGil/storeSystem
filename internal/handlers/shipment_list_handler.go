package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"storeSystem/internal/models"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetAllShipmentLists(w http.ResponseWriter, r *http.Request) {
	lists, err := h.shipmentListStore.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения листов отгрузки")
		return
	}

	respondWithJSON(w, http.StatusOK, lists)
}

func (h *Handlers) GetShipmentListByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id листа отгрузки")
		return
	}

	list, err := h.shipmentListStore.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, list)
}

func (h *Handlers) CreateShipmentList(w http.ResponseWriter, r *http.Request) {
	var input models.CreateShipmentListInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректно отправленные данные")
		return
	}

	if strings.TrimSpace(input.Article) == "" {
		respondWithError(w, http.StatusBadRequest, "article должен присутствовать")
		return
	}

	list, err := h.shipmentListStore.Create(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, list)
}

func (h *Handlers) AddShipmentFromFile(line string, userID int) (*models.ShipmentList, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("пустая строка")
	}

	params := strings.Split(line, ",")
	if len(params) < 4 {
		return nil, fmt.Errorf("некорректная строка: %s", line)
	}

	var input models.CreateShipmentListInput

	input.ShipmentId, _ = strconv.Atoi(params[0])
	input.CustomerId, _ = strconv.Atoi(params[1])
	input.Article = params[2]
	input.ExpectedAmount, _ = strconv.Atoi(params[3])

	input.CreatedBy = userID
	input.UpdatedBy = &userID

	list, err := h.shipmentListStore.Create(input)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (h *Handlers) UploadShipmentList(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Ошибка получения файла:", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	fmt.Printf("Загружен файл: %+v\n", handler.Filename)

	reader := bufio.NewReader(file)
	firstLine := true

	claims, _ := GetUserClaimsFromContext(r.Context())

	var result []*models.ShipmentList

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if firstLine {
			firstLine = false
			continue
		}

		fmt.Printf("%s len line %d", line, len(line))

		list, err := h.AddShipmentFromFile(line, claims.UserID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		result = append(result, list)
	}

	respondWithJSON(w, http.StatusOK, result)
}
