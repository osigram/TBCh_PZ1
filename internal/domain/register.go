package domain

import (
	"crypto/ecdsa"
	"errors"
)

func NewRegisterTransaction(key *ecdsa.PrivateKey) (*Transaction, error) {
	return newTransaction(RegisterTransaction, []byte{}, key)
}

func verifyPreviousRegisterTransactions(blockchain Blockchain, client Client) error {
	key := ecdsa.PublicKey(client)
	for _, block := range blockchain {
		if key.X.Int64() == block.Transaction.Client.X.Int64() && key.Y.Int64() == block.Transaction.Client.Y.Int64() {
			return errors.New("client is already registered")
		}
	}

	return nil
}
