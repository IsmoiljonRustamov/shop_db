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

type Product struct {
	ID           int64
	CategoryID   int64
	CategoryName string
	Name         string
	Price        float64
	ImageUrl     string
	CreatedAt    time.Time
	Images       []*ProductImage
}

type ProductImage struct {
	ID             int64
	ImageUrl       string
	SequenceNumber int32
}

type GetProductsParams struct {
	Limit  int32
	Page   int32
	Search string
}

type GetProductsResponse struct {
	Products []*Product
	count    int32
}

func (m *DBManager) CreateProduct(product *Product) (int64, error) {
	var productID int64
	query := `
		INSERT INTO products (
			category_id,
			name,
			price,
			image_url
		) VALUES ($1, $2, $3, $4)
		RETURNING id
		`

	row := m.db.QueryRow(
		query,
		product.CategoryID,
		product.Name,
		product.Price,
		product.ImageUrl,
	)

	err := row.Scan(&productID)
	if err != nil {
		return 0, err
	}

	queryInsertImage := `
		INSERT INTO product_image (
			product_id,
			image_url,
			sequnce_number
		) VALUES ($1,$2,$3)	
		`
	for _, image := range product.Images {
		_, err := m.db.Exec(
			queryInsertImage,
			productID,
			image.ImageUrl,
			image.SequenceNumber,
		)
		if err != nil {
			return 0, nil
		}
	}
	return productID, nil
}

func (m *DBManager) Get(id int64) (*Product, error) {
	var product Product

	product.Images = make([]*ProductImage, 0)

	query := `
		Select
		p.id,
		p.category.id,
		c.name,
		p.name,
		p.price,
		p.image_url,
		p.created_at
	FROM products p
	INNER JOIN categories c on c.id = p.category_id
	WHERE P.id=$1`

	row := m.db.QueryRow(query, id)

	err := row.Scan(
		&product.ID,
		&product.CategoryID,
		&product.CategoryName,
		&product.Name,
		&product.Price,
		&product.ImageUrl,
		&product.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	queryImage := `
		SELECT 
			id,
			image_url,
			sequence_number
		FROM product_images
		WHERE product_id = $1`

	rows, err := m.db.Query(queryImage, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var image ProductImage

		err := row.Scan(
			&image.ID,
			&image.ImageUrl,
			&image.SequenceNumber,
		)
		if err != nil {
			return nil, err
		}
		product.Images = append(product.Images, &image)
	}

	return &product, nil

}

func (m *DBManager) GetAllProducts(params *GetProductsParams) (*GetProductsResponse, error) {
	var res GetProductsResponse

	res.Products = make([]*Product, 0)

	filter := ""
	if params.Search != "" {
		filter = fmt.Sprintf("Where name ilike '%s'",
			"%"+params.Search+"%")
	}
	query := `SELECT
		p.id,
		p.category_id
		c.name,
		p.price,
		p.image_url,
		p.created_at
	FROM products p
	INNER JOIN categories c on c.id = p.category_if
    ` + filter + `
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2 
	`

	offset := (params.Page - 1) * params.Limit
	rows, err := m.db.Query(query, params.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Product

		err := rows.Scan(
			&p.ID,
			&p.CategoryID,
			&p.CategoryName,
			&p.Name,
			&p.Price,
			&p.ImageUrl,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		res.Products = append(res.Products, &p)
	}
	return &res, nil
}

func (m *DBManager) UpdateProduct(product *Product) error {
	query := `
		UPDATE products SET
			category_id = $1,
			name = $2,
			price = $3,
			image_ur = $4,
			WHERE id=$5`

	res, err := m.db.Exec(
		query,
		product.CategoryID,
		product.Name,
		product.Price,
		product.ImageUrl,
		product.ID,
	)

	if err != nil {
		return err
	}
	rowsCount, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return sql.ErrNoRows
	}

	queryDeleteImages := `DELETE FROM product_images WHERE product_id=$1`
	_, err = m.db.Exec(queryDeleteImages, product.ID)
	if err != nil {
		return err
	}

	queryInsertImage := `
		INSERT INTO product_images (
			product_id,
			image_url,
			sequence_number
			) VALUES ($1,$2,$3)
	`
	for _, image := range product.Images {
		_, err := m.db.Exec(
			queryInsertImage,
			product.ID,
			image.ImageUrl,
			image.SequenceNumber,
		)
		if err != nil {
			return err
		}
	}
	return nil

}

func (m *DBManager) DeleteProduct(id int64) error {
	queryDeleteImages := `DELETE FROM product_images WHERE 
		product_id=$1`
	_, err := m.db.Exec(queryDeleteImages, id)
	if err != nil {
		return err
	}
	queryDelete := `DELETE FROM products WHERE id=$1`
	res, err := m.db.Exec(queryDelete, id)
	if err != nil {
		return err
	}
	rowsCount, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsCount == 0 {
		return sql.ErrNoRows
	}
	return nil

}
