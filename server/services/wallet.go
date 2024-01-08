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

func CreateNewWallet(db *sql.DB, hm *HashMap, username string) (*Wallet, error) {
	hm.Mu.Lock()
	defer hm.Mu.Unlock()

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
	_, err = db.Exec(`
		INSERT INTO "wallets"("public_key", "owner")
		VALUES($1, $2)
	`, address.Hex(), username)
	if err != nil {
		return nil, err
	}

	hm.Value[address.Hex()] = true
	return &Wallet{
		PublicKey:  address.Hex(),
		PrivateKey: privateKeyHex,
	}, nil
}

func FetchWalletAddresses(db *sql.DB, name string) ([]string, error) {
	result, err := db.Query(`SELECT "public_key" FROM "wallets" WHERE "owner" = $1;`, name)
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
