package database

import (
	"database/sql"
	"fmt"
	"storeSystem/internal/models"

	"github.com/jmoiron/sqlx"
)

type CounterpartyStore struct {
	db *sqlx.DB
}

func NewCounterpartyStore(db *sqlx.DB) *CounterpartyStore {
	return &CounterpartyStore{db: db}
}

func (s *CounterpartyStore) GetAll() ([]models.Counterparty, error) {
	var counterparties []models.Counterparty
	query := `SELECT counterparty_id, full_name, legal_form, inn, kpp, ogrn, legal_address, bank_name, bik, bank_account_number, contact_person, phone, role_id FROM counterparty;`

	err := s.db.Select(&counterparties, query)
	fmt.Println(err, counterparties)

	if err != nil {
		return nil, err
	}
	return counterparties, nil
}

func (s *CounterpartyStore) GetByID(id int) (*models.Counterparty, error) {
	var counterparty models.Counterparty
	query := `SELECT * FROM counterparty where counterparty_id=$1;`

	err := s.db.Get(&counterparty, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("counterparty with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}
	return &counterparty, nil
}

func (s *CounterpartyStore) Create(input models.CreateCounterpartyInput) (*models.Counterparty, error) {
	var counterparty models.Counterparty

	query := `
	INSERT INTO counterparty (full_name, legal_form, inn, kpp, ogrn, legal_address, bank_name, bik, bank_account_number, contact_person, phone, role_id)
	VALUES ($1, $2, $3, $4, $5 ,$6, $7, $8, $9, $10, $11, $12)
	returning counterparty_id, full_name, legal_form, inn, kpp, ogrn, legal_address, bank_name, bik, bank_account_number, contact_person, phone, role_id;`

	err := s.db.QueryRowx(query, input.FullName, input.LegalForm, input.INN, input.KPP, input.OGRN, input.LegalAddress, input.BankName, input.BIK, input.BankAccountNumber, input.ContactPerson, input.Phone, input.RoleId).StructScan(&counterparty)

	if err != nil {
		return nil, err
	}
	return &counterparty, nil
}

//func (s *CounterpartyStore) Update(counterparty_id int, input models.UpdateCounterpartyInput) (*models.Counterparty, error) {
//	counterparty, err := s.GetByID(counterparty_id)
//	if err != nil {
//		return nil, err
//	}
//
//	query := `UPDATE item SET delivery_list_id = $1, supplier_id = $2, name= $3, article = $4, updated_at=$5
//	WHERE item_id = $6 returning item_id, delivery_list_id, supplier_id, name, article, updated_at;`
//
//	var updatedItem models.Item
//
//	err = s.db.QueryRowx(query, item.DeliveryListID, item.SupplierID, item.Name, item.Article, item.UpdatedAt, item.ID).StructScan(&updatedItem)
//	if err != nil {
//		return nil, err
//	}
//	return &updatedItem, nil
//}
