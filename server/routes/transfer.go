package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"bcDashboard/services"
)

type body struct {
	PrivateKey string  `json:"privateKey"`
	PublicKey  string  `json:"publicKey"`
	To         string  `json:"to"`
	Amount     float64 `json:"amount"`
	Id         int     `json:"id"`
}

type bodyToken struct {
	PrivateKey   string `json:"privateKey"`
	PublicKey    string `json:"publicKey"`
	To           string `json:"to"`
	Amount       int64  `json:"amount"`
	TokenAddress string `json:"tokenAddress"`
	Id           int    `json:"id"`
}

func TransferEthereum(db *sql.DB, hashmaps map[int]*services.HashMap) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := services.IsAuthenticated(r)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		var data body
		if err := decoder.Decode(&data); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		result := db.QueryRow(`SELECT "network_url" FROM "networks" WHERE id = $1`, data.Id)
		var network_url string
		if scanErr := result.Scan(&network_url); scanErr != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(scanErr.Error()))
			return
		}
		hashmaps[data.Id].Mu.RLock()
		client := hashmaps[data.Id].Client
		if data.PrivateKey[:2] == "0x" {
			data.PrivateKey = data.PrivateKey[2:]
		}
		txHash, err := services.TransferEthereum(client, data.PublicKey, data.To, data.PrivateKey, data.Amount)
		hashmaps[data.Id].Mu.RUnlock()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(txHash))
	})
}

func TransferTokens(db *sql.DB, hashmaps map[int]*services.HashMap) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := services.IsAuthenticated(r)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		var data bodyToken
		if err := decoder.Decode(&data); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		result := db.QueryRow(`SELECT "network_url" FROM "networks" WHERE id = $1`, data.Id)
		var network_url string
		if scanErr := result.Scan(&network_url); scanErr != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(scanErr.Error()))
			return
		}
		client := hashmaps[data.Id].Client
		if data.PrivateKey[:2] == "0x" {
			data.PrivateKey = data.PrivateKey[2:]
		}
		txHash, err := services.TransferTokens(client, data.PublicKey, data.To, data.TokenAddress, data.Amount, data.PrivateKey)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(txHash))
	})
}
