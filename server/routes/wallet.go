package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"regexp"

	"bcDashboard/services"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi/v5"
)

type response struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

func CreateWallet(db *sql.DB, hm *services.HashMap) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := services.IsAuthenticated(r)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte(err.Error()))
			return
		}
		wallet, err := services.CreateNewWallet(db, hm, name)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		data := response{
			Private: wallet.PrivateKey,
			Public:  wallet.PublicKey,
		}
		fmt.Println(data)
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
		addresses, err := services.FetchWalletAddresses(db, name)
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

func FetchBalance(client *ethclient.Client, db *sql.DB) http.HandlerFunc {
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
