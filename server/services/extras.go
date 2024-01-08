package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type HashMap struct {
	Mu    sync.RWMutex
	Value map[string]bool
}

func processBlock(client *ethclient.Client, blockNumber *big.Int, db *sql.DB, hm *HashMap) error {
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	ch := make(chan error, block.Transactions().Len())
	createdAt := time.Unix(int64(block.Time()), 0)
	for _, tx := range block.Transactions() {
		wg.Add(1)
		go processTransaction(client, tx, blockNumber, db, &wg, ch, hm, createdAt)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for err := range ch {
		if err != nil {
			return err
		}
	}
	return nil
}

func InitDb(dbAddress string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbAddress)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "users" (
			username TEXT PRIMARY KEY,
			password TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "check" (
			id SERIAL PRIMARY KEY,
			block_number SERIAL NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "wallets" (
			public_key TEXT PRIMARY KEY,
			owner TEXT REFERENCES "users" ("username")
		);
	`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "transactions" (
			id SERIAL PRIMARY KEY,
			to_address TEXT REFERENCES "wallets" ("public_key"),
			from_address TEXT NOT NULL,
			value TEXT NOT NULL,
			tokens TEXT,
			contract_id SERIAL REFERENCES "contracts" ("id"),
			block NUMERIC NOT NULL,
			created_at TIMESTAMPTZ,
			hash TEXT NOT NULL
			network_id SERIAL REFERENCES "networks" ("networks")
		);
	`)
	if err != nil {
		return nil, err
	}
	return db, err
}

func AddWalletAddresses(db *sql.DB) (*HashMap, error) {
	hm := HashMap{Value: map[string]bool{}}
	addressRows, err := db.Query(`SELECT "public_key" FROM "wallets"`)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	for addressRows.Next() {
		var address string
		scanErr := addressRows.Scan(&address)
		if scanErr != nil {
			fmt.Println(scanErr.Error())
			return nil, scanErr
		}
		hm.Value[address] = true
	}
	return &hm, nil
}

func BackgroundTask(client *ethclient.Client, db *sql.DB, hm *HashMap) {
	if len(hm.Value) == 0 {
		return
	}

	// latest block
	latest, err := client.BlockNumber(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return
	}

	// recent check
	recentRow := db.QueryRow(`SELECT "block_number" FROM "check" ORDER BY "created_at" DESC LIMIT 1;`)
	var recent uint64
	if err = recentRow.Scan(&recent); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = processBlock(client, big.NewInt(int64(latest)), db, hm)
			if err != nil {
				fmt.Println(err.Error())
				log.Fatal(err)
			}
			_, err = db.Exec(`
				INSERT INTO "check"("block_number")
				VALUES($1);
			`, latest)
			if err != nil {
				fmt.Println(err.Error())
				log.Fatal(err)
			}
			return
		}
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	// Concurrently process blocks
	var wg sync.WaitGroup
	ch := make(chan error, latest-recent)
	for i := recent + 1; i < latest+1; i++ {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()
			ch <- processBlock(client, big.NewInt(i), db, hm)
		}(int64(i))
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for val := range ch {
		if val != nil {
			fmt.Println(val.Error())
			log.Fatal(err)
		}
	}
	_, err = db.Exec(`
		INSERT INTO "check"("block_number")
		VALUES($1);
	`, latest)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
}
