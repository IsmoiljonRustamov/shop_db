package main

import (
	"database/sql"
	"fmt"
	"log"

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
	connStr := fmt.Sprintf("user = %s password = %s host = %s port = %d dbname = %s sslmode=disable", PostgresUser, PostgresPassword, PostgresHost, PostgresPort, PostgresDatabase)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to Connection:%v", err)
	}
	DBManagaer := NewDBManagaer(db)

	// td1, err := DBManagaer.CreateCustomer(&Customer{
	// 	first_name:   "Burxon",
	// 	last_name:    "Nuriddinov",
	// 	phone_number: "+999897717170",
	// 	gender:       1,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to Create Customer: %v", err)
	// }
	// fmt.Println(td1)

	customers,err := DBManagaer.GetAll(&GetAllParam{
		limit: 10,
		page: 2,
	})
	if err != nil {
		log.Fatalf("Failed to GetAll: %v", err)
	}
	for _,val := range  customers.customers {
		fmt.Println(*val) 	
	}




}
