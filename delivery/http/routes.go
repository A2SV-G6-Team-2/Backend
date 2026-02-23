package http

import (
	"net/http"
	"strings"
)

func RegisterDebtRoutes(mux *http.ServeMux, handler *DebtHandler) {
	if mux == nil || handler == nil {
		return
	}

	mux.HandleFunc("/debts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateDebt(w, r)
		case http.MethodGet:
			handler.ListDebts(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/debts/upcoming", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.ListUpcomingDebts(w, r)
	})

	mux.HandleFunc("/debts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/pay") {
			if r.Method != http.MethodPatch {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			handler.MarkDebtPaid(w, r)
			return
		}

		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.UpdateDebt(w, r)
	})
}
