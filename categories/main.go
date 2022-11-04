package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

var (
	PostgresUser = "postgres"
	PostgresPassword = "12345"
	PostgresHost = "localhost"
	PostgresPort = 5432
	PostgresDatabase = "shop_db"
)

func main() {

	connCtr := fmt.Sprintf("user = %s password = %s host = %s port = %d dbname = %s sslmode = disable", PostgresUser,
	PostgresPassword,PostgresHost,PostgresPort,PostgresDatabase)

	db, err := sql.Open("postgres", connCtr)
	if err != nil {
		log.Fatalf("Failed to open connection1: %v", err)
	}
	DBManager := NewDBManager(db)

	// td1, err := DBManager.CreateCategories(Category{
	// 	name: "zohid saidov",
	// 	image_url: "url_2",
	// })

	// if err != nil {
	// 	log.Fatalf("Failed to Create categories: %v", err)
	// }
	// fmt.Println(td)
	
	// td2, err = DBManager.Get(1)
	// if err != nil {
	// 	log.Fatalf("failed to Get categories: %v",err)
	// }
	// fmt.Println(td)

	// td3,err := DBManager.Update(&Category{
	// 	id: 12,
	// 	image_url: "url_1",
	// 	name: "Ismoiljon",
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to Updated: %v",err)
	// }
	// fmt.Println(td3)

	// err = DBManager.Delete(10)
	// if err != nil {
	// 	log.Fatalf("Failed to deleted: %v", err)
	// }

	td4, err := DBManager.GetAll()
	if err != nil {
		log.Fatalf("Failed to getall: %v", err)
	}
	for a, v := range td4 {
		fmt.Println("key: ",a, "value: ", v)
	}

}