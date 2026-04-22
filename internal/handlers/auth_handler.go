package handlers

import (
	"encoding/json"
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

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Неправильное имя пользователя или пароль")
		return
	}

	token, err := auth.GenerateToken(worker.ID, worker.RoleId, worker.Username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true на https
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id":       worker.ID,
		"username": worker.Username,
		"role_id":  worker.RoleId,
	})
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true на https
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "успешный выход",
	})
}

func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserClaimsFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "нет данных пользователя")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"id":       claims.UserID,
		"username": claims.Username,
		"role_id":  claims.RoleID,
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
