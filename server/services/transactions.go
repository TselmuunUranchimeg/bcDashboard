package services

import (
	"bcDashboard/token"
	"context"
	"database/sql"
	"errors"
	"math"
	"math/big"
	"strings"
	"sync"
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

func processTransaction(client *ethclient.Client, tx *types.Transaction, blockNumber *big.Int, db *sql.DB, wg *sync.WaitGroup, ch chan error, hm *HashMap, createdAt time.Time) {
	hm.Mu.RLock()
	defer hm.Mu.RUnlock()
	defer wg.Done()
	from, err := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
	if err != nil {
		ch <- err
		return
	}
	to := tx.To() // Is the contract address if it was a transfer event
	if to == nil {
		return
	}
	bytecode, err := getBytecode(client, *to, blockNumber)
	if err != nil {
		ch <- err
		return
	}
	value := getReadableValue(tx.Value())
	realBlock := blockNumber.String()
	// Contract address
	if len(bytecode) > 0 {
		contractAbi, err := abi.JSON(strings.NewReader(token.TokenMetaData.ABI))
		if err != nil {
			ch <- err
			return
		}
		transferHash := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Hex()
		receipt, err := getTransactionReceipt(client, tx.Hash())
		if err != nil {
			ch <- err
			return
		}
		for _, receiptLog := range receipt.Logs {
			if len(receiptLog.Topics) < 3 {
				continue
			}
			if receiptLog.Topics[0].Hex() != transferHash {
				return
			}
			if common.HexToAddress(receiptLog.Topics[1].Hex()).Hex() != from.Hex() {
				return
			}
			if !hm.Value[common.HexToAddress(receiptLog.Topics[2].Hex()).Hex()] {
				return
			}
			obj, err := contractAbi.Unpack("Transfer", receiptLog.Data)
			if err != nil {
				ch <- err
				return
			}
			tokens := getReadableValue(obj[0].(*big.Int))
			dbTx, err := db.Begin()
			if err != nil {
				ch <- err
				return
			}
			_, err = dbTx.Exec(`
				INSERT INTO "transactions"("from_address", "to_address", "value", "hash", "tokens", "contract", "block", "created_at")
				VALUES($1, $2, $3, $4, $5, $6, $7, $8)
			`, from.Hex(), common.HexToAddress(receiptLog.Topics[2].Hex()).Hex(), value, tx.Hash().Hex(), tokens, to.Hex(), realBlock, createdAt)
			if err != nil {
				ch <- err
				dbTx.Rollback()
				return
			}
			err = dbTx.Commit()
			if err != nil {
				ch <- err
				return
			}
		}
	} else {
		if !hm.Value[to.Hex()] {
			return
		}
		dbTx, err := db.Begin()
		if err != nil {
			ch <- err
			return
		}
		_, err = dbTx.Exec(`
			INSERT INTO "transactions"("to_address", "from_address", "value", "hash", "block", "created_at")
			VALUES($1, $2, $3, $4, $5, $6);
		`, to.Hex(), from.Hex(), value, tx.Hash().Hex(), realBlock, createdAt)
		if err != nil {
			ch <- err
			return
		}
		err = dbTx.Commit()
		if err != nil {
			ch <- err
			return
		}
	}
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
