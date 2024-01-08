package routes

import (
	"encoding/json"
	"net/http"

	"bcDashboard/services"

	"github.com/ethereum/go-ethereum/ethclient"
)

type body struct {
	PrivateKey string  `json:"privateKey"`
	PublicKey  string  `json:"publicKey"`
	To         string  `json:"to"`
	Amount     float64 `json:"amount"`
}

type bodyToken struct {
	PrivateKey   string `json:"privateKey"`
	PublicKey    string `json:"publicKey"`
	To           string `json:"to"`
	Amount       int64  `json:"amount"`
	TokenAddress string `json:"tokenAddress"`
}

func TransferEthereum(client *ethclient.Client) http.HandlerFunc {
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
		if data.PrivateKey[:2] == "0x" {
			data.PrivateKey = data.PrivateKey[2:]
		}
		txHash, err := services.TransferEthereum(client, data.PublicKey, data.To, data.PrivateKey, data.Amount)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(txHash))
	})
}

func TransferTokens(client *ethclient.Client) http.HandlerFunc {
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
