package services

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

func TestGetTransactions(t *testing.T) {
	dbAddress := "user=postgres password=tselmuun100 dbname=bcDashboard port=5432 host=localhost sslmode=disable"
	db, err := InitDb(dbAddress)
	if err != nil {
		t.Error(err)
	}
	result, err := GetTransactions(db, "0x9b3FDBBE4e112e3925b934c34698F2A695dFA43c")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(*result))
	for _, val := range *result {
		fmt.Printf("%+v\n", val)
	}
}

func TestProcessTransaction(t *testing.T) {
	client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/oX7PDYCyymtGX9QkcejJPNX_dfSnStSO")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	dbAddress := "user=postgres password=tselmuun100 dbname=bcDashboard port=5432 host=localhost sslmode=disable"
	db, err := InitDb(dbAddress)
	if err != nil {
		t.Error(err)
	}
	tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash("0xad25278729bf75c75d7bcf02d1b55fbe2c0c669d5545d5db17595b96177f433e"))
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan error, 1)
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
	blockNumber := big.NewInt(5004774)
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	createdAt := time.Unix(int64(block.Time()), 0)
	start := time.Now()
	go func() {
		defer wg.Done()
		ch <- processTransaction(client, db, tx, blockNumber, &hm, createdAt)
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

func BenchmarkProcessTransaction(b *testing.B) {
	client, err := ethclient.Dial("wss://eth-sepolia.g.alchemy.com/v2/oX7PDYCyymtGX9QkcejJPNX_dfSnStSO")
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()
	dbAddress := "user=postgres password=tselmuun100 dbname=bcDashboard port=5432 host=localhost sslmode=disable"
	db, err := InitDb(dbAddress)
	if err != nil {
		b.Fatal(err)
	}
	hashes := []string{"0xad25278729bf75c75d7bcf02d1b55fbe2c0c669d5545d5db17595b96177f433e", "0x310c780e85ac4e13d233f7a3f185499336f94801129963fac47d2f08925463ab"}
	blockNumbers := []*big.Int{big.NewInt(5004774), big.NewInt(5005042)}
	var wg sync.WaitGroup
	ch := make(chan error, 2)
	hm := HashMap{
		Value: map[string]bool{
			"0x9b3FDBBE4e112e3925b934c34698F2A695dFA43c": true,
		},
		Id: 1,
	}
	result, err := db.Query(`SELECT "contract" FROM "contracts" WHERE "network" = 1;`)
	if err != nil {
		b.Fatal(err)
	}
	for result.Next() {
		var contract string
		if scanErr := result.Scan(&contract); scanErr != nil {
			b.Fatal(scanErr)
		}
		hm.Contracts[contract] = true
	}
	for i := 0; i < 2; i++ {
		tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(hashes[i]))
		if err != nil {
			b.Fatal(err)
		}
		wg.Add(1)
		block, err := client.BlockByNumber(context.Background(), blockNumbers[i])
		if err != nil {
			b.Fatal(err)
		}
		ch <- processTransaction(client, db, tx, blockNumbers[i], &hm, time.Unix(int64(block.Time()), 0))
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for val := range ch {
		fmt.Println(val)
	}
}
