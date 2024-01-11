package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"bcDashboard/services"
)

type contractBody struct {
	Address   string `json:"address"`
	NetworkId int    `json:"networkId"`
}

func RegisterContract(db *sql.DB, hashmaps map[int]*services.HashMap) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := services.IsAuthenticated(r)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("not authenticated"))
			return
		}
		if name != os.Getenv("ADMIN_NAME") {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("not admin"))
			return
		}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var data contractBody
		err = decoder.Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("can't decode"))
			return
		}
		hm := hashmaps[data.NetworkId]
		_, err = db.Query(
			`SELECT "network_url" FROM "networks" WHERE "id" = $1 AND "network_url" = $2;`,
			data.NetworkId,
			hm.Url,
		)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("invalid network id"))
			return
		}
		if err = services.RegisterContract(db, data.Address, hm); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("Successfully registered the contract"))
	})
}
