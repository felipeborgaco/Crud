package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Product struct {
	Id          int
	Name        string
	Description string
	Price       float32
	Stock       int
}

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://paulo:123456@postgres/crud?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}
func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	u := Product{}

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Println("server failed to handle", err)
		return
	}

	_, err = db.Exec("INSERT INTO product (name, description,price, stock) VALUES ($1,$2,$3,$4)", u.Name, u.Description, u.Price, u.Stock)

	if err != nil {
		fmt.Println("server failed to handle", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func Read(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM product")
	if err != nil {
		fmt.Println("server failed to handle", err)
		return
	}
	defer rows.Close()

	data := make([]Product, 0)
	for rows.Next() {
		Product := Product{}
		err := rows.Scan(&Product.Id, &Product.Name, &Product.Description, &Product.Price, &Product.Stock)
		if err != nil {
			fmt.Println("server failed to handle", err)
		}
		data = append(data, Product)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("server failed to handle", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	up := Product{}
	err := json.NewDecoder(r.Body).Decode(&up)
	if err != nil {
		fmt.Println("server failed to handle", err)
		return
	}

	row := db.QueryRow("SELECT * FROM product WHERE id=$1", id)

	u := Product{}
	err = row.Scan(&u.Id, &u.Name, &u.Description, &u.Price, &u.Stock)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if up.Name != "" {
		u.Name = up.Name
	}

	if up.Description != "" {
		u.Description = up.Description
	}

	if up.Price != 0 {
		u.Price = up.Price
	}

	if up.Stock != 0 {
		u.Stock = up.Stock
	}

	_, err = db.Exec("UPDATE product SET name=$1, description=$2, price=$3, stock=$4 WHERE id=$5;", u.Name, u.Description, u.Price, u.Stock, u.Id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}
func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	_, err := db.Exec("DELETE FROM product WHERE id=$1;", id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Println("server starting")
	http.HandleFunc("/product/read", Read)
	http.HandleFunc("/product/create", Create)
	http.HandleFunc("/product/delete", Delete)
	http.HandleFunc("/product/update", Update)
	http.ListenAndServe(":8080", nil)
}
