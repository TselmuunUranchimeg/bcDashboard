package services

import (
	"crypto/ecdsa"
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	PublicKey, PrivateKey string
}

func CreateNewWallet(db *sql.DB, username string, networkId int) (*Wallet, error) {
	// Create private key
	pk, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	privateKeyBytes := crypto.FromECDSA(pk)
	privateKeyHex := hexutil.Encode(privateKeyBytes)

	// Public key
	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("can't convert")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Save the new wallet address
	dbTx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	_, err = dbTx.Exec(`
		INSERT INTO "wallets"("public_key", "owner", "network_id")
		VALUES($1, $2, $3)
	`, address.Hex(), username, networkId)
	if err != nil {
		dbTx.Rollback()
		return nil, err
	}
	if err = dbTx.Commit(); err != nil {
		dbTx.Rollback()
		return nil, err
	}
	return &Wallet{
		PublicKey:  address.Hex(),
		PrivateKey: privateKeyHex,
	}, nil
}

func FetchWalletAddresses(db *sql.DB, name, networkId string) ([]string, error) {
	result, err := db.Query(`SELECT "public_key" FROM "wallets" WHERE "owner" = $1 AND "network_id" = $2;`, name, networkId)
	if err != nil {
		return nil, err
	}
	addresses := []string{}
	for result.Next() {
		var address string
		scanErr := result.Scan(&address)
		if scanErr != nil {
			return nil, scanErr
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}
