package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
}

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://scs:123@postgres/crud?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS product (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL UNIQUE,
        description TEXT,
        price DECIMAL(10, 2) NOT NULL,
        stock INT NOT NULL
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database.")
}

func Create(w http.ResponseWriter, r *http.Request) {
	var p Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if p.Price < 0 {
		http.Error(w, "Price cannot be negative", http.StatusBadRequest)
		return
	}
	if p.Stock < 0 {
		http.Error(w, "Stock cannot be negative", http.StatusBadRequest)
		return
	}
	var existingId int
	err = db.QueryRow("SELECT id FROM product WHERE name = $1", p.Name).Scan(&existingId)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Failed to check product name uniqueness", http.StatusInternalServerError)
		return
	}
	if existingId != 0 {
		http.Error(w, "Product with the same name already exists", http.StatusBadRequest)
		return
	}
	err = db.QueryRow(
		"INSERT INTO product (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id",
		p.Name, p.Description, p.Price, p.Stock).Scan(&p.Id)
	if err != nil {
		http.Error(w, "Failed to insert product", http.StatusInternalServerError)
		return
	}

	// Retornando o produto criado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func Read(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr, exists := params["productId"]

	if !exists {
		rows, err := db.Query("SELECT * FROM product")
		if err != nil {
			http.Error(w, "Failed to query products", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		products := make([]Product, 0)
		for rows.Next() {
			var p Product
			err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.Stock)
			if err != nil {
				http.Error(w, "Failed to scan product", http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, "Failed to iterate over rows", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(products); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var p Product
	err = db.QueryRow("SELECT * FROM product WHERE id = $1", id).Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.Stock)
	if err == sql.ErrNoRows {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to query product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["productId"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var p Product
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if p.Price < 0 {
		http.Error(w, "Price cannot be negative", http.StatusBadRequest)
		return
	}
	if p.Stock < 0 {
		http.Error(w, "Stock cannot be negative", http.StatusBadRequest)
		return
	}

	var existingId int
	err = db.QueryRow("SELECT id FROM product WHERE name = $1 AND id != $2", p.Name, id).Scan(&existingId)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Failed to check product name uniqueness", http.StatusInternalServerError)
		return
	}
	if existingId != 0 {
		http.Error(w, "Product with the same name already exists", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE product SET name=$1, description=$2, price=$3, stock=$4 WHERE id=$5",
		p.Name, p.Description, p.Price, p.Stock, id)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	p.Id = id // Certifica-se de que o ID do produto atualizado seja retornado corretamente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["productId"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var p Product
	err = db.QueryRow("SELECT * FROM product WHERE id = $1", id).Scan(&p.Id, &p.Name, &p.Description, &p.Price, &p.Stock)
	if err == sql.ErrNoRows {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to query product", http.StatusInternalServerError)
		return
	}

	// Deletar o produto
	_, err = db.Exec("DELETE FROM product WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	fmt.Println("Server starting")

	r := mux.NewRouter()

	r.HandleFunc("/products", Create).Methods("POST")
	r.HandleFunc("/products", Read).Methods("GET")
	r.HandleFunc("/products/{productId}", Read).Methods("GET")
	r.HandleFunc("/products/{productId}", Update).Methods("PUT")
	r.HandleFunc("/products/{productId}", Delete).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}
