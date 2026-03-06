package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"storeSystem/internal/auth"
	"storeSystem/internal/models"
	"storeSystem/internal/validation"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "неправильные поля")
		return
	}

	worker, err := h.workerStore.GetByUsername(input.Username)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Имя пользователя не найдено")
		return
	}

	err = auth.CheckPassword(worker.PasswordHash, input.Password)
	fmt.Println(err)
	fmt.Println(worker.PasswordHash)
	fmt.Println(input.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Неправильное имя пользователя или пароль")
		return
	}

	role := ""

	switch worker.RoleId {
	case 1:
		role = "ADMIN"
	case 2:
		role = "WAREHOUSE_WORKER"
	default:
		role = "UNKNOWN"
	}

	token, err := auth.GenerateToken(worker.ID, role)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"username": worker.Username,
		"token":    token,
		"role_id":  worker.RoleId,
	})
}

func (h *Handlers) CreateWorker(w http.ResponseWriter, r *http.Request) {
	var input models.CreateWorkerInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	err := validation.Validate.Struct(input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.workerStore.GetByUsername(input.Username)
	if err == nil {
		respondWithError(w, http.StatusConflict, "user with this username already exists")
		return
	}

	worker, err := h.workerStore.Create(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, worker)
}
