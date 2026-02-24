package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"storeSystem/internal/models"

	"storeSystem/internal/auth"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	worker, err := h.workerStore.GetByUsername(input.Username)
	if err != nil {
		http.Error(w, "unauthorized getbyusername", http.StatusUnauthorized)
		return
	}

	err = auth.CheckPassword(worker.PasswordHash, input.Password)
	fmt.Println(err)
	fmt.Println(worker.PasswordHash)
	fmt.Println(input.Password)

	if err != nil {
		http.Error(w, "unauthorized passwordhash", http.StatusUnauthorized)
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
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *Handlers) CreateWorker(w http.ResponseWriter, r *http.Request) {
	var input models.CreateWorkerInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	worker, err := h.workerStore.Create(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, worker)
}
