package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func createHashmap(db *sql.DB, id int, url string) (*HashMap, error) {
	hm := HashMap{
		Value:     map[string]bool{},
		Id:        id,
		Contracts: map[string]bool{},
		Url:       url,
	}
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	hm.Client = client
	addressRows, err := db.Query(`SELECT "public_key" FROM "wallets" WHERE "network_id" = $1;`, id)
	if err != nil {
		return nil, err
	}
	for addressRows.Next() {
		var address string
		if err = addressRows.Scan(&address); err != nil {
			return nil, err
		}
		hm.Value[address] = true
	}
	contractRows, err := db.Query(`SELECT "contract" FROM "contracts" WHERE "network" = $1;`, id)
	if err != nil {
		return nil, err
	}
	for contractRows.Next() {
		var contract string
		if scanErr := contractRows.Scan(&contract); scanErr != nil {
			return nil, scanErr
		}
		hm.Contracts[contract] = true
	}
	return &hm, nil
}

func addCheck(db *sql.DB, latest uint64, networkId int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO "check"("block_number", "network") VALUES($1, $2);`, latest, networkId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func processBlock(client *ethclient.Client, db *sql.DB, block *types.Block, hm *HashMap) error {
	length := block.Transactions().Len()
	ch := make(chan error, length)
	var wg sync.WaitGroup
	wg.Add(length)
	for _, tx := range block.Transactions() {
		go func(tx *types.Transaction) {
			defer wg.Done()
			ch <- processTransaction(client, db, tx, block.Number(), hm, time.Unix(int64(block.Time()), 0))
		}(tx)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for val := range ch {
		if val != nil {
			return val
		}
	}
	return nil
}

func backgroundClientTask(db *sql.DB, hm *HashMap) error {
	latestBlock, err := hm.Client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	l := latestBlock.Number().Uint64()
	recentCheck := db.QueryRow(`SELECT "block_number" FROM "check" WHERE "network" = $1 ORDER BY "created_at" DESC LIMIT 1`, hm.Id)
	var recent uint64
	if err := recentCheck.Scan(&recent); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err = addCheck(db, l, hm.Id); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	difference := l - recent
	if difference == 0 {
		fmt.Printf("There is no new block on network %d\n", hm.Id)
		return nil
	}
	var wg sync.WaitGroup
	wg.Add(int(difference))
	ch := make(chan error, difference)
	fmt.Printf("Network %d will be processing %d blocks: ", hm.Id, difference)
	for i := recent + 1; i < l+1; i++ {
		if i == l {
			fmt.Printf("%d\n", i)
		} else {
			fmt.Printf("%d ", i)
		}
	}
	for i := recent + 1; i < l+1; i++ {
		go func(i int64) {
			defer wg.Done()
			block, err := hm.Client.BlockByNumber(context.Background(), big.NewInt(i))
			if err != nil {
				ch <- err
			}
			ch <- processBlock(hm.Client, db, block, hm)
		}(int64(i))
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for val := range ch {
		if val != nil {
			return val
		}
	}
	if err = addCheck(db, l, hm.Id); err != nil {
		return err
	}
	return nil
}
