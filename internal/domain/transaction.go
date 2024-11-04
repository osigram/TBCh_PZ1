package domain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"time"
)

type TransactionType int

const (
	RegisterTransaction TransactionType = iota
	SendTransaction
)

type Transaction struct {
	Timestamp time.Time
	Type      TransactionType
	Client    Client
	Data      []byte
	Signature []byte
}

func newTransaction(transactionType TransactionType, data []byte, key *ecdsa.PrivateKey) (*Transaction, error) {
	t := Transaction{
		Timestamp: time.Now(),
		Client:    Client(key.PublicKey),
		Type:      transactionType,
		Data:      data,
		Signature: nil,
	}

	hash, err := t.Hash()

	sign, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		return nil, err
	}
	t.Signature = sign

	return &t, err
}

func (t Transaction) VerifySignature(key ecdsa.PublicKey) error {
	hash, err := t.Hash()
	if err != nil {
		return err
	}

	if ecdsa.VerifyASN1(&key, hash[:], t.Signature) {
		return nil
	}

	return errors.New("wrong signature")
}

func (t Transaction) Hash() ([32]byte, error) {
	t.Signature = nil

	tBytes, err := json.Marshal(t)
	if err != nil {
		return [32]byte{}, err
	}

	hash := sha256.Sum256(tBytes)

	return hash, err
}

func Account(blockchain Blockchain, client Client) (int64, error) {
	var account int64

	key := ecdsa.PublicKey(client)

	for i := len(blockchain) - 1; i >= 0; i-- {
		block := &blockchain[i]
		if key.X.Int64() == block.Transaction.Client.X.Int64() && key.Y.Int64() == block.Transaction.Client.Y.Int64() && block.Transaction.Type == RegisterTransaction {
			account += 100
			break
		}
		if block.Transaction.Type != SendTransaction {
			continue
		}
		data, err := NewSendTransactionDataFromJSON(block.Transaction.Data)
		if err != nil {
			return 0, err
		}

		if key.X.Int64() == block.Transaction.Client.X.Int64() && key.Y.Int64() == block.Transaction.Client.Y.Int64() {
			account -= data.Amount
		}
		if key.X.Int64() == data.To.X.Int64() && key.Y.Int64() == data.To.Y.Int64() {
			account += data.Amount
		}
	}

	return account, nil
}

func verifyPreviousTransactions(blockchain Blockchain, transaction Transaction) error {
	if len(blockchain) == 0 {
		return nil
	}

	if transaction.Type == RegisterTransaction {
		return verifyPreviousRegisterTransactions(blockchain, transaction.Client)
	}

	if transaction.Type == SendTransaction {
		return verifyPreviousSendTransactions(blockchain, transaction)
	}

	return errors.New("unknown transaction type")
}
