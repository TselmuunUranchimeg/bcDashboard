package services

import (
	"database/sql"
	"fmt"
	"testing"

	"bcDashboard/services"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestRegisterContract(t *testing.T) {
	db, err := sql.Open("postgres", "sslmode=disable user=postgres password=tselmuun100 dbname=bcDashboard host=localhost port=5432")
	if err != nil {
		t.Fatal(err)
	}
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/oX7PDYCyymtGX9QkcejJPNX_dfSnStSO")
	if err != nil {
		t.Fatal(err)
	}
	hm := services.HashMap{
		Client:    client,
		Contracts: map[string]bool{},
		Value:     map[string]bool{},
		Id:        1,
	}
	if err = services.RegisterContract(db, "0xc45df6dd7a0D3AeD8133E9977d5c8a608f9f2d13", &hm); err != nil {
		t.Fatal(err)
	}
	fmt.Println(hm.Contracts)
}
