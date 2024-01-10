package services

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

func TestProcessBlock(t *testing.T) {
	client, err := ethclient.Dial("https://binance.llamarpc.com")
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("postgres", "user=postgres password=tselmuun100 dbname=bcDashboard port=5432 host=localhost sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	addressResult, err := db.Query(`SELECT "public_key" FROM "wallets";`)
	if err != nil {
		t.Fatal(err)
	}
	hm := HashMap{Value: map[string]bool{}}
	for addressResult.Next() {
		var address string
		if scanErr := addressResult.Scan(&address); scanErr != nil {
			t.Fatal(scanErr)
		}
		hm.Value[address] = true
	}
	result, err := db.Query(`SELECT "contract" FROM "contracts" WHERE "network" = 1;`)
	if err != nil {
		t.Fatal(err)
	}
	for result.Next() {
		var contract string
		if scanErr := result.Scan(&contract); scanErr != nil {
			t.Fatal(scanErr)
		}
		hm.Contracts[contract] = true
	}
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan error, 1)
	start := time.Now()
	go func() {
		defer wg.Done()
		ch <- processBlock(client, db, block, &hm)
	}()
	go func() {
		wg.Wait()
		close(ch)
	}()
	for val := range ch {
		if val != nil {
			t.Fatal(val)
		}
	}
	fmt.Println(time.Since(start).String())
}
