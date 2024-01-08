package services

import (
	"bcDashboard/token"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func getNecessaryData(client *ethclient.Client, privateKey string) (privateKeyECDSA *ecdsa.PrivateKey, nonce uint64, gasPrice, chainId *big.Int, err error) {
	privateKeyECDSA, err = crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, 0, nil, nil, err
	}
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, 0, nil, nil, errors.New("can't get public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err = client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, 0, nil, nil, err
	}
	gasPrice, err = client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, 0, nil, nil, err
	}
	chainId, err = client.ChainID(context.Background())
	if err != nil {
		return nil, 0, nil, nil, err
	}
	return privateKeyECDSA, nonce, gasPrice, chainId, nil
}

func floatToBigInt(value float64) *big.Int {
	v := new(big.Float)
	v.SetFloat64(value)
	v.Mul(v, big.NewFloat(math.Pow10(18)))
	result := new(big.Int)
	v.Int(result)
	return result
}

func TransferEthereum(client *ethclient.Client, from, to, privateKey string, value float64) (string, error) {
	// Check balance
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(from), nil)
	if err != nil {
		return "", err
	}
	amount := floatToBigInt(value)
	if balance.Cmp(amount) == -1 {
		return "", errors.New("you don't have enough funds for this transaction")
	}

	privateKeyECDSA, nonce, gasPrice, chainId, err := getNecessaryData(client, privateKey)
	if err != nil {
		return "", err
	}
	toAddress := common.HexToAddress(to)
	math.Pow10(18)
	tx := types.NewTransaction(nonce, toAddress, amount, 300000, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), privateKeyECDSA)
	if err != nil {
		return "", err
	}
	if err = client.SendTransaction(context.Background(), signedTx); err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}

func TransferTokens(client *ethclient.Client, from, to, contract string, value int64, privateKey string) (string, error) {
	// Check token balance
	tokenAddress := common.HexToAddress(contract)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		return "", err
	}
	balance, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(from))
	if err != nil {
		return "", err
	}
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return "", err
	}
	valueString := fmt.Sprintf("%d", value)
	for i := 0; i < int(decimals); i++ {
		valueString = valueString + "0"
	}
	amount := new(big.Int)
	amount, ok := amount.SetString(valueString, 10)
	if !ok {
		return "", errors.New("can't convert")
	}
	if balance.Cmp(amount) == -1 {
		return "", errors.New("you don't have enough funds for this transaction")
	}

	// Get token amount
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	// Get other necessary details
	privateKeyECDSA, nonce, gasPrice, chainId, err := getNecessaryData(client, privateKey)
	if err != nil {
		return "", err
	}
	toAddress := common.HexToAddress(to)
	transferSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferSignature)
	methodId := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	var data []byte
	data = append(data, methodId...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// Create, sign and send the transaction
	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), 300000, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), privateKeyECDSA)
	if err != nil {
		return "", err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}
