package database

import (
	"database/sql"
	"fmt"
	"net/http"
	"storeSystem/internal/auth"
	"storeSystem/internal/models"

	"github.com/jmoiron/sqlx"
)

type WorkerStore struct {
	db *sqlx.DB
}

func (s *WorkerStore) GetByUsername(username string) (*models.Worker, error) {
	var worker models.Worker

	query := `
	SELECT worker_id, full_name, username, password_hash, role_id
	FROM worker
	WHERE username=$1
	`

	err := s.db.Get(&worker, query, username)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("worker not found")
	}

	if err != nil {
		return nil, err
	}

	return &worker, nil
}

func (s *WorkerStore) Create(worker models.CreateWorkerInput) (*models.Worker, error) {
	hash, err := auth.HashPassword(worker.Password)
	if err != nil {
		return nil, err
	}

	query := `
	INSERT INTO worker (full_name, username, password_hash, role_id)
	VALUES ($1,$2,$3,$4)
	RETURNING worker_id, full_name, username, role_id
	`

	var createdWorker models.Worker

	err = s.db.QueryRowx(
		query,
		worker.FullName,
		worker.Username,
		hash,
		worker.RoleId,
	).StructScan(&createdWorker)

	if err != nil {
		return nil, err
	}

	return &createdWorker, nil
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := auth.GetUserFromContext(r)

		if user.Role != "ADMIN" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NewWorkerStore(db *sqlx.DB) *WorkerStore {
	return &WorkerStore{db: db}
}
