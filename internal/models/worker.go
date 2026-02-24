package models

type Worker struct {
	ID           int    `json:"worker_id" db:"worker_id"`
	FullName     string `json:"full_name" db:"full_name"`
	Username     string `json:"username" db:"username"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	RoleId       int    `json:"role_id" db:"role_id"`
}

type CreateWorkerInput struct {
	FullName string `json:"full_name" db:"full_name"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	RoleId   int    `json:"role_id" db:"role_id"`
}
