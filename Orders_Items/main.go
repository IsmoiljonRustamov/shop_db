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
		log.Fatalf("failed to open connection: %v", err)
	}

	m := NewDBManager(db)

	// _, err = m.CreateOrders(&Orders{
	// 	customer_id:  1,
	// 	total_amount: 100,
	// 	address:       "Chilonzor",
	// 	Items: []*Orders_items{
	// 		{
	// 			product_id:   2,
	// 			order_id: 2,
	// 			count:        2,
	// 			total_price:  100.00,
	// 			product_name: "Jacket",
	// 		},
	// 		{
	// 			product_id:   2,
	// 			count:        3,
	// 			total_price:  140.00,
	// 			product_name: "Shoes",
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	log.Fatalf("failed to create Items: %v", err)
	// }

	// order, err := m.GetOrders(id)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(order)

	// orders, err := m.GetAllOrders(&GetOrdersParams{
	// 	limit: 10,
	// 	page: 1,
	// 	search: "Jacket",
	// })
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(orders)

	m.DeleteOrders(17)
}
