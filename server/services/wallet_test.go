package services

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestCreateNewWallet(t *testing.T) {
	db, err := InitDb("user=postgres password=tselmuun100 dbname=bcDashboard port=5432 host=localhost sslmode=disable")
	if err != nil {
		t.Error(err)
	}
	hm := &HashMap{Value: map[string]bool{}}
	wallet, err := CreateNewWallet(db, hm, "Tselmuun")
	if err != nil {
		t.Error(err)
	}
	result := db.QueryRow(`SELECT * FROM "wallets" WHERE "public_key" = $1;`, wallet.PublicKey)
	var (
		pbKey, owner string
	)
	err = result.Scan(&pbKey, &owner)
	if err != nil {
		t.Error(err)
	}
	if pbKey != wallet.PublicKey && owner != "Tselmuun" {
		t.Errorf("Got %s instead of %s\n", pbKey, wallet.PublicKey)
	}
}
