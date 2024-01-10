package services

import (
	"fmt"
	"testing"

	"bcDashboard/services"
)

func TestCreateNewWallet(t *testing.T) {
	db, err := services.InitDb("user=postgres password=tselmuun100 dbname=bcDashboard port=5432 host=localhost sslmode=disable")
	if err != nil {
		t.Error(err)
	}
	hm := &services.HashMap{
		Value: map[string]bool{},
		Id:    2,
	}
	wallet, err := services.CreateNewWallet(db, "Tselmuun", hm.Id)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(wallet)
	result := db.QueryRow(`SELECT * FROM "wallets" WHERE "public_key" = $1;`, wallet.PublicKey)
	var (
		pbKey, owner string
		networkId    int
	)
	err = result.Scan(&pbKey, &owner, &networkId)
	if err != nil {
		t.Error(err)
	}
	if pbKey != wallet.PublicKey && owner != "Tselmuun" && networkId != 2 {
		t.Errorf("Got %s instead of %s\n", pbKey, wallet.PublicKey)
	}
}
