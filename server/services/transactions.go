package services

import (
	"bcDashboard/token"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getReadableValue(val *big.Int) string {
	f, _ := val.Float64()
	num := new(big.Float).Quo(big.NewFloat(f), big.NewFloat(math.Pow10(18)))
	return num.String()
}

func getBytecode(client *ethclient.Client, address common.Address, blockNumber *big.Int) ([]byte, error) {
	bytecode, err := client.CodeAt(context.Background(), address, blockNumber)
	if err != nil {
		if strings.Contains(err.Error(), "Your app has exceeded its compute units") {
			i := 1
			for err != nil {
				time.Sleep(time.Second * time.Duration(i))
				code, err := client.CodeAt(context.Background(), address, blockNumber)
				if err == nil {
					return code, err
				}
				i += 1
			}
		}
	}
	return bytecode, err
}

func getTransactionReceipt(client *ethclient.Client, hash common.Hash) (*types.Receipt, error) {
	receipt, err := client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		if strings.Contains(err.Error(), "Your app has exceeded its compute units") {
			i := 1
			for err != nil {
				time.Sleep(time.Duration(i) * time.Second)
				r, err := client.TransactionReceipt(context.Background(), hash)
				if err == nil {
					return r, err
				}
				i += 1
			}
		}
	}
	return receipt, err
}

func processTransaction(client *ethclient.Client, db *sql.DB, tx *types.Transaction, blockNumber *big.Int, hm *HashMap, createdAt time.Time) error {
	hm.Mu.RLock()
	defer hm.Mu.RUnlock()
	from, err := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
	if err != nil {
		return err
	}
	to := tx.To()
	if to == nil {
		return nil
	}
	bytecode, err := getBytecode(client, *to, blockNumber)
	if err != nil {
		return err
	}
	value := tx.Value()
	// Contract address
	if len(bytecode) > 0 {
		if len(hm.Contracts) == 0 {
			return nil
		}
		if !hm.Contracts[to.Hex()] {
			return nil
		}
		contractAbi, err := abi.JSON(strings.NewReader(token.TokenMetaData.ABI))
		if err != nil {
			return err
		}
		transferHash := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Hex()
		receipt, err := getTransactionReceipt(client, tx.Hash())
		if err != nil {
			return err
		}
		for _, receiptLog := range receipt.Logs {
			if len(receiptLog.Topics) < 3 {
				continue
			}
			if receiptLog.Topics[0].Hex() != transferHash {
				return nil
			}
			if common.HexToAddress(receiptLog.Topics[1].Hex()).Hex() != from.Hex() {
				return nil
			}
			recipient := common.HexToAddress(receiptLog.Topics[2].Hex()).Hex()
			if !hm.Value[recipient] {
				return nil
			}
			obj, err := contractAbi.Unpack("Transfer", receiptLog.Data)
			if err != nil {
				return err
			}
			tokens := getReadableValue(obj[0].(*big.Int))
			transaction, err := db.Begin()
			if err != nil {
				return err
			}
			_, err = transaction.Exec(`
				INSERT INTO "transactions"("to_address", "from_address", "value", "valueWei", "block", "hash", "network_id", "created_at", "tokens", "contract")
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
			`, recipient, from.Hex(), getReadableValue(value), value.String(), blockNumber.Uint64(), tx.Hash().Hex(), hm.Id, createdAt, tokens, to.Hex())
			if err != nil {
				transaction.Rollback()
				return err
			}
			if err = transaction.Commit(); err != nil {
				transaction.Rollback()
				return err
			}
		}
		return nil
	}
	if !hm.Value[to.Hex()] {
		return nil
	}
	transaction, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = transaction.Exec(`
		INSERT INTO "transactions"("to_address", "from_address", "value", "valueWei", "block", "hash", "network_id", "created_at")
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`, to.Hex(), from.Hex(), getReadableValue(value), value.String(), blockNumber.Uint64(), tx.Hash().Hex(), hm.Id, createdAt)
	if err != nil {
		transaction.Rollback()
		fmt.Printf("Recipient is %s. This message is coming from the regular Ethereum transfer section\n", to.Hex())
		return err
	}
	if err = transaction.Commit(); err != nil {
		transaction.Rollback()
		return err
	}
	return nil
}

type transaction struct {
	From        string `json:"from"`
	Contract    string `json:"contract"`
	Value       string `json:"value"`
	Hash        string `json:"hash"`
	Tokens      string `json:"tokens"`
	Created_at  string `json:"createdAt"`
	BlockNumber string `json:"block"`
}

func GetTransactions(db *sql.DB, address string) (*[]transaction, error) {
	result, err := db.Query(`
			SELECT "from_address", "value", "tokens", "contract", "created_at", "hash", "block"
			FROM "transactions" WHERE "to_address" = $1 ORDER BY "created_at" DESC LIMIT 20;
		`, address)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &[]transaction{}, nil
		}
		return nil, err
	}
	defer result.Close()
	transactions := []transaction{}
	for result.Next() {
		var (
			from, value, hash, block string
			created_at               time.Time
			tokens                   sql.NullString
			contract                 sql.NullString
		)
		scanErr := result.Scan(&from, &value, &tokens, &contract, &created_at, &hash, &block)
		if scanErr != nil {
			return nil, scanErr
		}
		transactions = append(transactions, transaction{
			From:        from,
			Contract:    contract.String,
			Value:       value,
			Tokens:      tokens.String,
			Created_at:  created_at.String(),
			Hash:        hash,
			BlockNumber: block,
		})
	}
	return &transactions, nil
}
