package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"bcDashboard/services"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
)

type response struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}
type createBody struct {
	Id int `json:"id"`
}

func CreateWallet(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := services.IsAuthenticated(r)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		var body createBody
		if err = decoder.Decode(&body); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		wallet, err := services.CreateNewWallet(db, name, body.Id)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			if err, ok := err.(*pq.Error); ok {
				if err.Code.Name() == "foreign_key_violation" {
					if strings.Contains(err.Error(), "owner") {
						w.Write([]byte("Owner doesn't exist"))
					} else {
						w.Write([]byte("Network is not registered"))
					}
					return
				}
			}
			w.Write([]byte(err.Error()))
			return
		}
		data := response{
			Private: wallet.PrivateKey,
			Public:  wallet.PublicKey,
		}
		res, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(res)
	})
}

func FetchAddresses(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := services.IsAuthenticated(r)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}
		id := chi.URLParam(r, "network")
		addresses, err := services.FetchWalletAddresses(db, name, id)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		data, err := json.Marshal(addresses)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(data)
	})
}

func FetchBalance(db *sql.DB, hashmaps map[int]*services.HashMap) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := services.IsAuthenticated(r)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}
		address := chi.URLParam(r, "address")
		reg := regexp.MustCompile("^0x[a-zA-Z0-9]{40}$")
		if !reg.MatchString(address) {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Not a valid address"))
			return
		}
		_, err = db.Query(`SELECT * FROM "wallets" WHERE "owner" = $1 AND "public_key" = $2`, name, address)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("this wallet is not associated with you"))
				return
			}
		}
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Invalid network id"))
			return
		}
		result := db.QueryRow(`SELECT "network_url" FROM "networks" WHERE "id" = $1`, id)
		var url string
		if err = result.Scan(&url); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		client := hashmaps[id].Client
		balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}
		fbalance := new(big.Float)
		fbalance.SetString(balance.String())
		ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18))).String()
		w.Write([]byte(ethValue))
	})
}
