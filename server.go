package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type item struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type server struct {
	*mux.Router
	shoppingItems []item
}

func NewServer() *server {
	s := &server{
		Router:        mux.NewRouter(),
		shoppingItems: []item{},
	}
	s.routes()
	return s
}

func (s *server) routes() {
	s.HandleFunc("/shopping-list", s.listShoppingItems()).Methods("GET")
	s.HandleFunc("/shopping-list", s.createShoppingItem()).Methods("POST")
	s.HandleFunc("/shopping-list/{id}", s.removeShoppingItem()).Methods("DELETE")
}

func (s *server) createShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i item
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		i.ID = uuid.New()
		s.shoppingItems = append(s.shoppingItems, i)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) listShoppingItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.shoppingItems); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) removeShoppingItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idstr, ok := mux.Vars(r)["id"]
		if !ok {
			http.Error(w, "ID parameter is missing", http.StatusBadRequest)
			return
		}
		id, err := uuid.Parse(idstr)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}
		found := false
		for i, item := range s.shoppingItems {
			if item.ID == id {
				s.shoppingItems = append(s.shoppingItems[:i], s.shoppingItems[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
