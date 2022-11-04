package main

import (
	"database/sql"
	_"fmt"
	"time"
)


type DBManager struct {
	db sql.DB
}

func NewDBManager(db *sql.DB) *DBManager {
	return &DBManager{db:*db}
}

type Category struct {
	id int
	name string
	image_url string
	created_at time.Time
}

type GettAllParam struct {
	limit int
	page int
	title string
	assignee string
}


func (d *DBManager) CreateCategories(ct Category) (*Category,error) {
	query := `
		INSERT INTO categories (
			name,
			image_url
		) VALUES ($1, $2)
		RETURNING id, name, image_url, created_at`

	row := d.db.QueryRow(query,
		ct.name,
		ct.image_url,
	)
	
	var res Category
	err := row.Scan(
		&res.id,
		&res.name,
		&res.image_url,
		&res.created_at,
	)
	if err != nil {
		return nil, err 
	}
	return &res,nil 

}

func (d *DBManager) Get(id int) (*Category, error) {
	var res Category

	query := `
		SELECT 
		id,
		name,
		image_url,
		created_at
		FROM categories where id = $1
		`
	row := d.db.QueryRow(query, id)

	err := row.Scan(
		&res.id,
		&res.name,
		&res.image_url,
		&res.created_at,
	)

	if err != nil {
		return nil, err 
	}
	return &res, nil

}

func (d *DBManager) Update(ct *Category) (*Category, error) {
	query := `UPDATE categories SET  
			image_url = $1,
			name = $2
			WHERE id = $3
			RETURNING id, name, image_url, created_at`

	row := d.db.QueryRow(
		query,
		ct.image_url,
		ct.name,
		ct.id,
	)		
	var res Category
	err := row.Scan(
		&res.id,
		&res.image_url,
		&res.name,
		&res.created_at,
	)
	if err != nil {
		return nil, err
	}
	return ct,nil
}

func (d *DBManager) Delete(id int) (error) {
	query := `DELETE FROM categories 
		WHERE id=$1`

	_,err  := d.db.Exec(query,id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBManager) GetAll() ([]*Category,error) {
	query := `
		SELECT
		id,
		name,
		image_url,
		created_at
		FROM categories`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err 
	}
	defer rows.Close()
	var res []*Category

	for rows.Next() {
		var r Category

		err := rows.Scan(
			&r.id,
			&r.name,
			&r.image_url,
			&r.created_at,

		)
		if err != nil {
			return nil, err 
		}
		res = append(res, &r)
	}
	return 	res,nil 
}

