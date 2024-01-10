package services

import (
	"bcDashboard/token"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
)

func RegisterContract(db *sql.DB, address string, hm *HashMap) error {
	start := time.Now()
	instance, err := token.NewToken(common.HexToAddress(address), hm.Client)
	if err != nil {
		return err
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		return err
	}
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return err
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		return err
	}
	hm.Mu.Lock()
	defer hm.Mu.Unlock()
	dbTx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = dbTx.Exec(`
			INSERT INTO "contracts"("contract", "decimals", "symbol", "name", "network")
			VALUES ($1, $2, $3, $4, $5);
		`, address, decimals, symbol, name, hm.Id)
	if err != nil {
		dbTx.Rollback()
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return errors.New("contract is already registered")
			}
		}
		return err
	}
	if err = dbTx.Commit(); err != nil {
		dbTx.Rollback()
		return err
	}
	hm.Contracts[address] = true
	fmt.Println(time.Since(start).String())
	return nil
}
