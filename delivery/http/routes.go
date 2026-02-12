package http

import (
	"net/http"
	"strings"
)

// Router holds handlers and serves routes
type Router struct {
	Expense  *ExpenseHandler
	Category *CategoryHandler
}

// NewRouter creates a new router with the given handlers
func NewRouter(expense *ExpenseHandler, category *CategoryHandler) *Router {
	return &Router{Expense: expense, Category: category}
}

// Handler returns the main HTTP handler with middleware applied
func (rt *Router) Handler() http.Handler {
	mux := http.NewServeMux()

	// Expenses: POST /expenses, GET /expenses, GET/PUT/DELETE /expenses/:id
	mux.HandleFunc("/expenses", rt.expenseBase)
	mux.HandleFunc("/expenses/", rt.expenseByID)

	// Categories: POST /categories, GET /categories, GET/PUT/DELETE /categories/:id
	mux.HandleFunc("/categories", rt.categoryBase)
	mux.HandleFunc("/categories/", rt.categoryByID)

	// API documentation (Swagger UI) at /api-docs
	ServeAPIDocs(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"service":"Personal Expense Tracker","expenses":"/expenses","categories":"/categories"}`))
	})

	return MockUserIDMiddleware(mux)
}

func (rt *Router) expenseBase(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/expenses" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodPost:
		rt.Expense.Create(w, r)
	case http.MethodGet:
		rt.Expense.List(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (rt *Router) expenseByID(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/expenses/" {
		http.NotFound(w, r)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/expenses/")
	if idx := strings.Index(id, "/"); idx >= 0 {
		id = id[:idx]
	}
	switch r.Method {
	case http.MethodGet:
		rt.Expense.GetByID(w, r, id)
	case http.MethodPut:
		rt.Expense.Update(w, r, id)
	case http.MethodDelete:
		rt.Expense.Delete(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (rt *Router) categoryBase(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/categories" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodPost:
		rt.Category.Create(w, r)
	case http.MethodGet:
		rt.Category.List(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (rt *Router) categoryByID(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/categories/" {
		http.NotFound(w, r)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/categories/")
	if idx := strings.Index(id, "/"); idx >= 0 {
		id = id[:idx]
	}
	switch r.Method {
	case http.MethodGet:
		rt.Category.GetByID(w, r, id)
	case http.MethodPut:
		rt.Category.Update(w, r, id)
	case http.MethodDelete:
		rt.Category.Delete(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
