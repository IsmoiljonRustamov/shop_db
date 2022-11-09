package main

import (
	"database/sql"
	"fmt"
	"time"
)

type DBManager struct {
	db *sql.DB
}

func NewDBManager(db *sql.DB) DBManager {
	return DBManager{db}
}

type Order struct {
	Id          int
	CustomerId  int
	TotalAmount float64
	CreatedAt   time.Time
	Address     string
	Items       []*OrderItems
}

type OrderItems struct {
	Id          int
	OrderId     int
	ProductId   int
	Count       int
	TotalPrice  float64
	ProductName string
}

type GetOrdersParams struct {
	Limit  int32
	Page   int32
	Search string
}

type GetOrdersRespone struct {
	Orders []*Order
	Count  int32
}

func (m *DBManager) CreateOrders(or *Order) (int, error) {
	tx, err := m.db.Begin()

	var orderID int

	query := `INSERT INTO orders (
		customer_id,
		total_amount,
		address
	)VALUES ($1,$2,$3)
		RETURNING id
	`

	row := tx.QueryRow(query,
		or.CustomerId,
		or.TotalAmount,
		or.CreatedAt,
	)
	err = row.Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	queryInsertItems := `
		INSERT INTO order_items (
			product_id,
			order_id,
			count,
			total_price,
			product_name
		) VALUES($1,$2,$3,$4,$5)
			
		`
	for _, item := range or.Items {
		_, err := tx.Exec(
			queryInsertItems,
			item.ProductId,
			orderID,
			item.Count,
			item.TotalPrice,
			item.ProductName,
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return orderID, nil
}

func (m *DBManager) GetOrder(id int) (*Order, error) {
	var order Order

	order.Items = make([]*OrderItems, 0)

	query := `
		SELECT 
		id,
		customer_id,
		total_amount,
		created_at,
		address
		FROM orders WHERE id = $1
		`
	row := m.db.QueryRow(query, id)

	err := row.Scan(
		&order.Id,
		&order.CustomerId,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.Address,
	)
	if err != nil {
		return nil, err
	}

	queryItems := `
		SELECT
			id,
			order_id,
			product_id,
			count,
			total_price,
			product_name
		FROM order_items
		WHERE order_id=$1`

	rows, err := m.db.Query(queryItems, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var orit OrderItems

		err := rows.Scan(
			&orit.Id,
			&orit.OrderId,
			&orit.ProductId,
			&orit.Count,
			&orit.TotalPrice,
			&orit.ProductName,
		)
		fmt.Println(orit)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, &orit)
	}
	return &order, nil
}

func (m *DBManager) GetAllOrders(params *GetOrdersParams) (*GetOrdersRespone, error) {
	var res GetOrdersRespone

	res.Orders = make([]*Order, 0)
	offset := (params.Page - 1) * params.Limit
	filter := ""
	if params.Search != "" {
		filter = fmt.Sprintf("WHERE customer_id = %s", params.Search)
	}

	query := `
	SELECT 
		o1.id,
		o1.customer_id,
		o1.address,
		o1.total_amount
	FROM orders o1  
	` + filter + `
	ORDER BY o1.created_at ASC LIMIT $1 OFFSET $2
	`
	rows, err := m.db.Query(query, params.Limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.Id,
			&order.CustomerId,
			&order.Address,
			&order.TotalAmount,
		)
		if err != nil {
			return nil, err
		}
		res.Orders = append(res.Orders, &order)
	}
	return &res, nil
}

func (m *DBManager) UpdateOrders(orders *Order) error {
	tx, err := m.db.Begin()
	query := `
		UPDATE orders SET total_amount = $1
			WHERE customer_id = $2`

	result, err := m.db.Exec(
		query,
		orders.TotalAmount,
		orders.CustomerId,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	res, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if res == 0 {
		tx.Rollback()
		return sql.ErrNoRows
	}
	query = `DELETE FROM order_items WHERE order_id=$1`
	_, err = m.db.Exec(query, orders.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	query = `
	 	INSERT INTO order_items (
			order_id,
			product_name,
			product_id,
			count,
			total_price,
		) VALUES($1,$2,$3,$4,$5)`
	for _, v := range orders.Items {
		_, err := m.db.Exec(
			query,
			v.OrderId,
			v.ProductName,
			v.ProductId,
			v.Count,
			v.TotalPrice,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (m *DBManager) DeleteOrder(order_id int) error {
	query := `DELETE FROM order_items WHERE order_id=$1`
	row, err := m.db.Exec(query, order_id)
	if err != nil {
		return err
	}
	res, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if res == 0 {
		return sql.ErrNoRows
	}

	queryDeleteOrders := `DELETE FROM orders WHERE id = $1`
	_, err = m.db.Exec(queryDeleteOrders, order_id)
	if err != nil {
		return err
	}
	return nil

}
