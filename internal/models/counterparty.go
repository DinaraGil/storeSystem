package models

type Counterparty struct {
	ID                int    `json:"counterparty_id" db:"counterparty_id"`
	FullName          string `json:"full_name" db:"full_name"`
	LegalForm         string `json:"legal_form" db:"legal_form"`
	INN               string `json:"inn" db:"inn"`
	KPP               string `json:"kpp" db:"kpp"`
	OGRN              string `json:"ogrn" db:"ogrn"`
	LegalAddress      string `json:"legal_address" db:"legal_address"`
	BankName          string `json:"bank_name" db:"bank_name"`
	BIK               string `json:"bik" db:"bik"`
	BankAccountNumber string `json:"bank_account_number" db:"bank_account_number"`
	ContactPerson     string `json:"contact_person" db:"contact_person"`
	Phone             string `json:"phone" db:"phone"`
	RoleId            int    `json:"role_id" db:"role_id"`
}

type CreateCounterpartyInput struct {
	FullName          string `json:"full_name" db:"full_name"`
	LegalForm         string `json:"legal_form" db:"legal_form"`
	INN               string `json:"inn" db:"inn"`
	KPP               string `json:"kpp" db:"kpp"`
	OGRN              string `json:"ogrn" db:"ogrn"`
	LegalAddress      string `json:"legal_address" db:"legal_address"`
	BankName          string `json:"bank_name" db:"bank_name"`
	BIK               string `json:"bik" db:"bik"`
	BankAccountNumber string `json:"bank_account_number" db:"bank_account_number"`
	ContactPerson     string `json:"contact_person" db:"contact_person"`
	Phone             string `json:"phone" db:"phone"`
	RoleId            int    `json:"role_id" db:"role_id"`
}

type UpdateCounterpartyInput struct {
	FullName          *string `json:"full_name" db:"full_name"`
	LegalForm         *string `json:"legal_form" db:"legal_form"`
	INN               *string `json:"inn" db:"inn"`
	KPP               *string `json:"kpp" db:"kpp"`
	OGRN              *string `json:"ogrn" db:"ogrn"`
	LegalAddress      *string `json:"legal_address" db:"legal_address"`
	BankName          *string `json:"bank_name" db:"bank_name"`
	BIK               *string `json:"bik" db:"bik"`
	BankAccountNumber *string `json:"bank_account_number" db:"bank_account_number"`
	ContactPerson     *string `json:"contact_person" db:"contact_person"`
	Phone             *string `json:"phone" db:"phone"`
	RoleId            *int    `json:"role_id" db:"role_id"`
}
