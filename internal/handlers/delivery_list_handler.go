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

func (h *Handlers) GetAllDeliveryLists(w http.ResponseWriter, r *http.Request) {
	delLists, err := h.deliveryListStore.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения листов поставок")
		return
	}
	respondWithJSON(w, http.StatusOK, delLists)
}

func (h *Handlers) GetDeliveryListByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный id листа поставки")
		return
	}
	delList, err := h.deliveryListStore.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, delList)
}

func (h *Handlers) CreateDeliveryList(w http.ResponseWriter, r *http.Request) {
	var input models.CreateDeliveryListInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректно отправленные данные")
		return
	}

	if strings.TrimSpace(input.Article) == "" {
		respondWithError(w, http.StatusBadRequest, "article должен присутствовать")
		return
	}

	delList, err := h.deliveryListStore.Create(input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, delList)
}

func (h *Handlers) AddFromFile(line string, userID int) (*models.DeliveryList, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("пустая строка")
	}

	params := strings.Split(line, ",")
	if len(params) < 4 {
		return nil, fmt.Errorf("некорректная строка: %s", line)
	}

	var input models.CreateDeliveryListInput

	input.DeliveryId, _ = strconv.Atoi(params[0])
	input.SupplierId, _ = strconv.Atoi(params[1])
	input.Article = params[2]
	input.ExpectedAmount, _ = strconv.Atoi(params[3])
	input.CreatedBy = userID
	input.UpdatedBy = &userID

	delList, err := h.deliveryListStore.Create(input)

	if err != nil {
		return nil, err
	}
	return delList, nil
}

func (h *Handlers) UploadDeliveryList(w http.ResponseWriter, r *http.Request) {
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

	var result []*models.DeliveryList

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

		delList, err := h.AddFromFile(line, claims.UserID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		result = append(result, delList)
	}

	respondWithJSON(w, http.StatusOK, result)

}
