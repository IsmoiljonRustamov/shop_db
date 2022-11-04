package main

import (
	"database/sql"
	"log"
	"fmt"
	_ "github.com/lib/pq"
)

var (
	PostgresUser     = "postgres"
	PostgresPassword = "12345"
	PostgresHost     = "localhost"
	PostgresPort     = 5432
	PostgresDatabase = "shop_db"
)

func main() {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		PostgresHost,
		PostgresPort,
		PostgresUser,
		PostgresPassword,
		PostgresDatabase,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open connection: %v",err)
	}

	m := NewDBManager(db)

	id, err := m.CreateProduct(&Product{
		CategoryID: 2,
		Name: "Iphone 13",
		Price: 900.0,
		ImageUrl: "test_url",
		Images: []*ProductImage{
			{
				ImageUrl: "test_url_1",
				SequenceNumber: 1,
			},
			{
				ImageUrl: "test_url_2",
				SequenceNumber: 2,
			},
			{
				ImageUrl: "test_url_3",
				SequenceNumber: 3,
			},
		},
	})
	if err != nil {
		log.Fatalf("failed to create product: %v" , err)
	}

	product , err := m.Get(id)
	if err != nil {
		log.Fatalf("failed to get product: %v" ,err)
	}
	fmt.Println(product)
}