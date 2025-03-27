package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatal("Error occurred iniatialising DB")
	}
	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products(
	ID int AUTO_INCREMENT, 
	Name VARCHAR(255) NOT NULL, 
	Quantity int, 
	Price float(10,7),
	PRIMARY KEY (ID)
	);`

	_, err := a.DB.Query(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) VALUES('%v', %v, %v)", name, quantity, price)
	a.DB.Exec(query)
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 20, 111)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

}

func checkStatusCode(t *testing.T, expectStatusCode int, actualStatusCode int) {
	if expectStatusCode != actualStatusCode {
		t.Errorf("Expected status %v, Received status %v", expectStatusCode, actualStatusCode)

	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	var product = []byte(`{"name": "Book", "quantity": 21, "price":200}`)

	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "Book" {
		t.Errorf("Expect name %v, Received name %v", "Book", m["name"])
	}
	if m["quantity"] != 21.0 {
		t.Errorf("Expect quantity %v, Received quantity %v", 21, m["quantity"])
	}
	if m["price"] != 200.0 {
		t.Errorf("Expect price %v, Received price %v", 200, m["price"])
	}

}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)

}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 10)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	var product = []byte(`{"name": "connector", "quantity": 1, "price":10}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response = sendRequest(req)
	if response.Code != 200 {
		t.Error(response.Code)
	}
	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["id"] != newValue["id"] {
		t.Errorf("Got oldValue %v, and newValue %v", oldValue["id"], newValue["id"])
	}
}
