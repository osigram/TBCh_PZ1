package domain

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
)

type SendTransactionData struct {
	To     Client
	Amount int64
}

func NewSendTransactionDataFromJSON(data []byte) (*SendTransactionData, error) {
	var result SendTransactionData
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func NewSendTransaction(key *ecdsa.PrivateKey, to Client, amount int64) (*Transaction, error) {
	sendData := SendTransactionData{
		To:     to,
		Amount: amount,
	}

	data, err := json.Marshal(sendData)
	if err != nil {
		return nil, err
	}

	return newTransaction(SendTransaction, data, key)
}

func verifyPreviousSendTransactions(blockchain Blockchain, transaction Transaction) error {
	if len(blockchain) == 0 {
		return errors.New("empty blockchain")
	}

	account, err := Account(blockchain, transaction.Client)
	if err != nil {
		return err
	}

	data, err := NewSendTransactionDataFromJSON(transaction.Data)
	if err != nil {
		return err
	}

	if data.Amount > account {
		return errors.New("insufficient funds")
	}

	return nil
}
