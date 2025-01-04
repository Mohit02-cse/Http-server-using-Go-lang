package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type item struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
}

type server struct {
	*mux.Router
	db *sql.DB
}

func NewServer(dsn string) (*server, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS Shopping_List"); err != nil {
		return nil, err
	}

	if _, err := db.Exec("USE Shopping_List"); err != nil {
		return nil, err
	}

	s := &server{
		Router: mux.NewRouter(),
		db:     db,
	}
	s.routes()
	return s, nil
}

func (s *server) routes() {
	s.HandleFunc("/shopping-list/{customer}", s.listShoppingItems()).Methods("GET")
	s.HandleFunc("/shopping-list/{customer}", s.createShoppingItem()).Methods("POST")
	s.HandleFunc("/shopping-list/{customer}", s.updateItemQuantity()).Methods("PUT")
	s.HandleFunc("/shopping-list/{customer}", s.removeShoppingItem()).Methods("DELETE")
}

func (s *server) createShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customer := mux.Vars(r)["customer"]
		if customer == "" {
			http.Error(w, "Customer name is required", http.StatusBadRequest)
			return
		}

		var i item
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if i.Name == "" {
			http.Error(w, "Item name is required", http.StatusBadRequest)
			return
		}

		tableName := fmt.Sprintf("`%s`", customer)
		createTableQuery := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				id CHAR(36) PRIMARY KEY,
				name VARCHAR(255) UNIQUE NOT NULL,
				quantity INT NOT NULL DEFAULT 1
			)`, tableName)

		if _, err := s.db.Exec(createTableQuery); err != nil {
			http.Error(w, "Failed to create table: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var existingQuantity int
		err := s.db.QueryRow(fmt.Sprintf("SELECT quantity FROM %s WHERE name = ?", tableName), i.Name).Scan(&existingQuantity)

		if err == nil {
			updateQuery := fmt.Sprintf("UPDATE %s SET quantity = quantity + ? WHERE name = ?", tableName)
			if _, err := s.db.Exec(updateQuery, 1, i.Name); err != nil {
				http.Error(w, "Failed to update item quantity: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			i.ID = uuid.New()
			insertQuery := fmt.Sprintf("INSERT INTO %s (id, name, quantity) VALUES (?, ?, ?)", tableName)
			if _, err := s.db.Exec(insertQuery, i.ID.String(), i.Name, i.Quantity); err != nil {
				http.Error(w, "Failed to insert item: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(i)
	}
}

func (s *server) updateItemQuantity() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customer := mux.Vars(r)["customer"]
		if customer == "" {
			http.Error(w, "Customer name is required", http.StatusBadRequest)
			return
		}

		var payload struct {
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if payload.Name == "" || payload.Quantity == 0 {
			http.Error(w, "Item name and valid quantity are required", http.StatusBadRequest)
			return
		}

		tableName := fmt.Sprintf("`%s`", customer)
		updateQuery := fmt.Sprintf("UPDATE %s SET quantity = quantity + ? WHERE name = ?", tableName)

		result, err := s.db.Exec(updateQuery, payload.Quantity, payload.Name)
		if err != nil {
			http.Error(w, "Failed to update item quantity: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *server) listShoppingItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customer := mux.Vars(r)["customer"]
		if customer == "" {
			http.Error(w, "Customer name is required", http.StatusBadRequest)
			return
		}

		query := fmt.Sprintf("SELECT id, name, quantity FROM `%s`", customer)
		rows, err := s.db.Query(query)
		if err != nil {
			http.Error(w, "Failed to fetch items: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []item
		for rows.Next() {
			var i item
			var idStr string
			if err := rows.Scan(&idStr, &i.Name, &i.Quantity); err != nil {
				http.Error(w, "Failed to parse items: "+err.Error(), http.StatusInternalServerError)
				return
			}
			i.ID, _ = uuid.Parse(idStr)
			items = append(items, i)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}
}

func (s *server) removeShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customer := mux.Vars(r)["customer"]
		if customer == "" {
			http.Error(w, "Customer name is required", http.StatusBadRequest)
			return
		}

		itemName := r.URL.Query().Get("name")
		if itemName == "" {
			http.Error(w, "Item name is required", http.StatusBadRequest)
			return
		}

		tableName := fmt.Sprintf("`%s`", customer)
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE name = ?", tableName)

		result, err := s.db.Exec(deleteQuery, itemName)
		if err != nil {
			http.Error(w, "Failed to delete item: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
