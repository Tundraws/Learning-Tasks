package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price, stock FROM products")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err.Error())
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Database error", err.Error())
			return
		}
		products = append(products, p)
	}

	respondWithJSON(w, http.StatusOK, products)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	var p Product
	err = db.QueryRow("SELECT id, name, price, stock FROM products WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock)

	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Product not found", "")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	defer r.Body.Close()

	result, err := db.Exec("INSERT INTO products(name, price, stock) VALUES(?, ?, ?)",
		p.Name, p.Price, p.Stock)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err.Error())
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err.Error())
		return
	}

	p.ID = int(id)
	respondWithJSON(w, http.StatusCreated, p)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	var p Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}
	defer r.Body.Close()

	p.ID = id
	res, err := db.Exec("UPDATE products SET name = ?, price = ?, stock = ? WHERE id = ?",
		p.Name, p.Price, p.Stock, p.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err.Error())
		return
	}

	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Product not found", "")
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID", err.Error())
		return
	}

	res, err := db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err.Error())
		return
	}

	if rowsAffected == 0 {
		respondWithError(w, http.StatusNotFound, "Product not found", "")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func respondWithError(w http.ResponseWriter, code int, error string, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: error, Message: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
