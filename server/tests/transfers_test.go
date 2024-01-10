package services

import (
	"fmt"
	"testing"

	"bcDashboard/services"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestTransferEthereum(t *testing.T) {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/6UlG8ZHeQPuCSXhB1fDgPfKB0M91iyIs")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	hash, err := services.TransferEthereum(
		client,
		"0x885A80eDE1e25Ef5A8917580f0aF8EAbDBE87f9F",
		"0xA2862D00525Bc367c476a4f1AE5C44d8B87b6DA2",
		"16437b6de0cd1d43472e39f7cbda993cc8b7d0ff9a2d8699347e1a77c8e45fde",
		1.55,
	)
	fmt.Println(hash, err)
	if err == nil {
		t.Fatal("something went wrong with the check")
	}
}

func TestTransferTokens(t *testing.T) {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/6UlG8ZHeQPuCSXhB1fDgPfKB0M91iyIs")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	txHash, err := services.TransferTokens(
		client,
		"0xA2862D00525Bc367c476a4f1AE5C44d8B87b6DA2", // from
		"0x885A80eDE1e25Ef5A8917580f0aF8EAbDBE87f9F", //to,
		"0xc45df6dd7a0D3AeD8133E9977d5c8a608f9f2d13", //contract
		43000,
		"15687ea8ccda2cf1c95ee622f55fa39741cc56f31fe74f052ef217e246c03a1e", //private key
	)
	fmt.Println(txHash, err)
	if err == nil {
		t.Fatal("something went wrong with the check")
	}
}
