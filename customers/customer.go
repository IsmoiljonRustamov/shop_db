package main

import (
	"database/sql"
	"time"
)




type Customer struct {
	id int
	first_name string
	last_name string 
	phone_number string
	gender int
	birth_date time.Time
	created_at time.Time
	updated_at time.Time
	deleted_at time.Time
}

type DBManager struct {
	db sql.DB
}

func NewDBManagaer(db *sql.DB) *DBManager {
	return &DBManager{db:*db}
} 

type GetAllParam struct {
	limit int
	page int
	search string
}

func (b *DBManager) CreateCustomer(cus *Customer) (*Customer,error) {
	query := `INSERT INTO customers (
		first_name,
		last_name,
		phone_number,
		gender,
		birth_date
		)VALUES($1,$2,$3,$4,$5)
		RETURNING id, first_name, last_name,phone_number,gender,birth_date,created_at`

	row := 	b.db.QueryRow(query,
		cus.first_name,
		cus.last_name,
		cus.phone_number,
		cus.gender,
		cus.birth_date,)
		
	var res Customer
	
	err := row.Scan(
		&res.id,
		&res.first_name,
		&res.last_name,
		&res.phone_number,
		&res.gender,
		&res.birth_date,
		&res.created_at,
	)
	if err != nil {
		return nil,err 
	}
	return &res, nil 
}

func(d *DBManager) Get(id int) (*Customer,error) {
	var res Customer

	query := `
	SELECT 
	id,
	first_name,
	last_name,
	phone_number,
	gender,
	birth_date,
	created_at
	FROM customers`

	row := d.db.QueryRow(query,id)

	err := 	row.Scan(
		&res.id,
		&res.first_name,
		&res.last_name,
		&res.phone_number,
		&res.gender,
		&res.birth_date,
		&res.created_at,
	)
	if err != nil {
		return nil ,err 
	}
	return &res, nil 
}

func (d *DBManager) UpdateCustomer(upt *Customer) (*Customer,error) {
	query := `UPDATE customers SET 
		updated_at=$1
		WHERE id=$2,
		RETURNING id,updated_at`

	row := d.db.QueryRow(query,
		upt.id,
		upt.created_at)	
	var res Customer
	err := row.Scan(
		&res.id,
		&res.updated_at,
	)	
	if err != nil {
		return nil, err 
	}
	return &res, nil 
}

func (d *DBManager) Deleted_at(delt *Customer) {
	query := `UPDATE customers SET deleted_at=$1 WHERE id=$2`

	d.db.Exec(
		query,
		delt.id,
		delt.deleted_at,
	)
}



