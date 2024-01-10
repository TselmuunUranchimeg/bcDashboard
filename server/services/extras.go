package services

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type HashMap struct {
	Mu               sync.RWMutex
	Value, Contracts map[string]bool
	Id               int
	Client           *ethclient.Client
	Url              string
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
		CREATE TABLE IF NOT EXISTS "networks" (
			id SERIAL PRIMARY KEY,
			network_url TEXT NOT NULL,
			name TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "contracts" (
			contract TEXT PRIMARY KEY,
			decimals REAL NOT NULL,
			symbol TEXT NOT NULL,
			name TEXT NOT NULL,
			network SERIAL REFERENCES "networks" ("id")
		);
	`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "check" (
			id SERIAL PRIMARY KEY,
			block_number NUMERIC NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			network SERIAL REFERENCES "networks" ("id")
		);
	`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "wallets" (
			public_key TEXT PRIMARY KEY,
			owner TEXT REFERENCES "users" ("username"),
			network_id SERIAL REFERENCES "networks" ("id")
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
			value NUMERIC NOT NULL,
			valueWei NUMERIC NOT NULL,
			tokens TEXT,
			contract TEXT REFERENCES "contracts" ("contract"),
			block NUMERIC NOT NULL,
			created_at TIMESTAMPTZ,
			hash TEXT NOT NULL,
			network_id SERIAL REFERENCES "networks" ("id")
		);
	`)
	if err != nil {
		return nil, err
	}
	return db, err
}

func InitNetworks(db *sql.DB) (map[int]string, error) {
	result, err := db.Query(`SELECT "id", "network_url" FROM "networks" ORDER BY "id";`)
	if err != nil {
		return nil, err
	}
	networks := map[int]string{}
	for result.Next() {
		var (
			id  int
			url string
		)
		if scanErr := result.Scan(&id, &url); scanErr != nil {
			return nil, scanErr
		}
		networks[id] = url
	}
	return networks, nil
}

func InitHashmaps(db *sql.DB, networks map[int]string) (map[int]*HashMap, error) {
	result := map[int]*HashMap{}
	for key, val := range networks {
		hm, err := createHashmap(db, key, val)
		if err != nil {
			return nil, err
		}
		result[key] = hm
	}
	return result, nil
}

func BackgroundTask(db *sql.DB, hm *HashMap, duration int) {
	for {
		if err := backgroundClientTask(db, hm); err != nil {
			fmt.Printf("%v\n", err.Error())
			if hm.Client == nil {
				i := 1
				for hm.Client != nil {
					time.Sleep(time.Duration(i) * time.Second)
					hm.Client, err = ethclient.Dial(hm.Url)
					if err != nil {
						i += 1
					}
				}
			}
		}
		time.Sleep(time.Second * time.Duration(duration))
	}
}
