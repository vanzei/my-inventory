package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "SELECT id, name, quantity, price from products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil

}

func (p *product) getProduct(db *sql.DB) error {
	query := "SELECT name, quantity, price FROM products WHERE id = ?"
	row := db.QueryRow(query, p.ID)
	err := row.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	// fmt.Printf("Scanned product: %+v\n", p)
	return nil

}

func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO products (name, quantity, price) VALUE ('%v', %v, %v)", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	query := "UPDATE products SET name=?, quantity=?, price=? WHERE id=?"
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price, p.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("No row existed")
	}
	return nil

}

func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM products WHERE id=%v", p.ID)
	_, err := db.Exec(query)
	return err
}
