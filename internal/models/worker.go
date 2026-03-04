package models

type Worker struct {
	ID           int    `json:"worker_id" db:"worker_id"`
	FullName     string `json:"full_name" db:"full_name"`
	Username     string `json:"username" db:"username"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	RoleId       int    `json:"role_id" db:"role_id"`
}

type CreateWorkerInput struct {
	FullName string `json:"full_name" db:"full_name" validate:"required"`
	Username string `json:"username" db:"username" validate:"required"`
	Password string `json:"password" db:"password" validate:"required"`
	RoleId   int    `json:"role_id" db:"role_id" validate:"required,min=1,max=2"`
}
