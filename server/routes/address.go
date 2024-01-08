package routes

import (
	"bcDashboard/services"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

func FetchDetails(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := services.IsAuthenticated(r)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Not authenticated"))
			return
		}
		address := chi.URLParam(r, "address")
		reg := regexp.MustCompile("^0x[a-zA-Z0-9]{40}$")
		if !reg.MatchString(address) {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Not a valid wallet address."))
			return
		}
		_, err = db.Query(`SELECT "public_key" FROM "wallets" WHERE "owner" = $1 AND "public_key" = $2;`, name, address)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("You don't have any wallets at the moment."))
				return
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		transactions, err := services.GetTransactions(db, address)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		data, err := json.Marshal(*transactions)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(data)
	})
}
