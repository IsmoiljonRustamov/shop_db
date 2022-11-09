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

	// _, err = m.CreateOrders(&Order{
	// 	CustomerId:  1,
	// 	TotalAmount: 100,
	// 	Address:       "Chilonzor",
	// 	Items: []*OrderItems{
	// 		{
	// 			ProductId:   2,
	// 			OrderId: 2,
	// 			Count:        2,
	// 			TotalPrice:  100.00,
	// 			ProductName: "Jacket",
	// 		},
	// 		{
	// 			ProductId:   2,
	// 			Count:        3,
	// 			TotalPrice:  140.00,
	// 			ProductName: "Shoes",
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	log.Fatalf("failed to create Items: %v", err)
	// }

	// order, err := m.GetOrder(67)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(order.Id)
	// fmt.Println(order.TotalAmount)
	// fmt.Println(order.Address)
	// fmt.Println(order.CreatedAt)
	// fmt.Println(order.CustomerId)
	// fmt.Println(order.Items)

	// orders, err := m.GetAllOrders(&GetOrdersParams{
	// 	Limit: 10,
	// 	Page: 1,
	// 	Search: "1",
	// })
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(orders)

	err = m.DeleteOrder(68)
	if err != nil {
		log.Fatalln(err)
	}

	// err = m.UpdateOrders(&Order{
	// 	Id: 68,
	// 	CustomerId: 1,
	// 	TotalAmount: 100,
	// })
	// if err != nil {
	// 	log.Fatalln(err) 
	// }


}
