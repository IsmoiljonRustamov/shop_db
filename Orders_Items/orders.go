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

type Orders struct {
	id           int
	customer_id  int
	total_amount float64
	created_at   time.Time
	address      string
	Items        []*Orders_items
}

type Orders_items struct {
	id           int
	order_id     int
	product_id   int
	count        int
	total_price  float64
	product_name string
}

type GetOrdersParams struct {
	limit  int32
	page   int32
	search string
}

type GetOrdersRespone struct {
	orders []*Orders_items
	count  int32
}

func (m *DBManager) CreateOrders(or *Orders) (int, error) {
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
		or.customer_id,
		or.total_amount,
		or.created_at,
	)
	err = row.Scan(&orderID)
	tx.Rollback()
	if err != nil {
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
			item.product_id,
			orderID,
			item.count,
			item.total_price,
			item.product_name,
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return orderID, nil
}

func (m *DBManager) GetOrders(id int) (*Orders, error) {
	var order Orders

	order.Items = make([]*Orders_items, 0)

	query := `
		SELECT 
		id,
		customer_id,
		total_amount,
		created_at,
		address
		from orders where id = $1
		`
	row := m.db.QueryRow(query, id)

	err := row.Scan(
		&order.id,
		&order.customer_id,
		&order.total_amount,
		&order.created_at,
		&order.address,
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
	rows.Close()
	for rows.Next() {
		var orit Orders_items

		err := row.Scan(
			&orit.id,
			&orit.order_id,
			&orit.product_id,
			&orit.count,
			&orit.total_price,
			&orit.product_name,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, &orit)
	}
	return &order, nil
}

func (m *DBManager) GetAllOrders(params *GetOrdersParams) (*GetOrdersRespone, error) {
	var res GetOrdersRespone

	res.orders = make([]*Orders_items, 0)
	offset := (params.page - 1) * params.limit
	filter := ""
	if params.search != "" {
		filter = fmt.Sprintf("WHERE product_name like '%s'",
			"%"+params.search+"%")
	}

	query := `
	SELECT 
		o2.id,
		o1.id,
		o2.product_name,
		o2.product_id,
		o2.count,
		o1.total_amount
	FROM orders o1 JOIN order_items o2 ON o2.order_id = o1.id
	` + filter + `
	ORDER BY o2.id asc LIMIT $1 OFFSET $2
	`
	rows, err := m.db.Query(query, params.limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var orders Orders_items
		err := rows.Scan(
			&orders.id,
			&orders.order_id,
			&orders.product_name,
			&orders.product_id,
			&orders.count,
			&orders.total_price,
		)
		if err != nil {
			return nil, err
		}
		res.orders = append(res.orders, &orders)
	}
	return &res, nil
}

func (m *DBManager) UpdateOrders(orders *Orders) error {
	tx, err := m.db.Begin()
	query := `
		UPDATE orders SET total_amount = $1
			WHERE customer_id = $2`

	result, err := m.db.Exec(
		query,
		orders.total_amount,
		orders.customer_id,
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
	_, err = m.db.Exec(query, orders.id)
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
		) VALUES($1,$2,$3,$4,$5`
	for _, v := range orders.Items {
		_, err := m.db.Exec(
			query,
			v.order_id,
			v.product_name,
			v.product_id,
			v.count,
			v.total_price,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (m *DBManager) DeleteOrders(order_id int) error {
	query := `DELETE FROM order_items WHERE order_id=$!`
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
