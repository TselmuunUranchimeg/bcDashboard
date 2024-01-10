package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type res struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func FetchNetworks(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := db.Query(`SELECT "id", "name" FROM "networks";`)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		networks := []res{}
		for result.Next() {
			var (
				id   int
				name string
			)
			if scanErr := result.Scan(&id, &name); scanErr != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(err.Error()))
				return
			}
			networks = append(networks, res{
				Id:   id,
				Name: name,
			})
		}
		data, err := json.Marshal(networks)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(data)
	})
}
