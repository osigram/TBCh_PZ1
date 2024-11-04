package domain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"github.com/google/uuid"
)

const BlockHashZerosNum = 8

type Block struct {
	ID          uuid.UUID
	PrevHash    []byte
	Hash        []byte
	Nonce       uint64
	Transaction Transaction
}

func NewBlock(blockchain *Blockchain, transaction Transaction) (*Block, error) {
	var prevHash []byte
	if blockchain == nil || len(*blockchain) == 0 {
		prevHash = nil
	} else {
		prevHash = (*blockchain)[len(*blockchain)-1].Hash
	}

	block := Block{
		ID:          uuid.New(),
		PrevHash:    prevHash,
		Hash:        nil,
		Nonce:       0,
		Transaction: transaction,
	}

	nonce, hash, err := generateHash(block)
	if err != nil {
		return nil, err
	}
	block.Nonce = nonce
	block.Hash = hash[:]

	if blockchain != nil {
		*blockchain = append(*blockchain, block)
	}

	return &block, nil
}

func NewBlockWithValidation(blockchain *Blockchain, transaction Transaction) (*Block, error) {
	block, err := NewBlock(blockchain, transaction)
	if err != nil {
		return nil, err
	}
	if err := block.Validate(*blockchain); err != nil {
		return nil, err
	}

	return block, nil
}

func (b Block) Validate(blockchain Blockchain) error {
	if err := b.Transaction.VerifySignature(ecdsa.PublicKey(b.Transaction.Client)); err != nil {
		return err
	}
	if err := verifyPreviousTransactions(blockchain, b.Transaction); err != nil {
		return err
	}

	return nil
}

func (b Block) CalculateHash() (hash [32]byte, err error) {
	b.Hash = nil

	bBytes, err := json.Marshal(b)
	if err != nil {
		return hash, err
	}
	hash = sha256.Sum256(bBytes)

	return hash, err
}

func generateHash(block Block) (nonce uint64, hash [32]byte, err error) {
	hash[0] = 0b11111111
	for hash[0]>>(8-BlockHashZerosNum) != 0 {
		block.Nonce++

		hash, err = block.CalculateHash()
		if err != nil {
			return nonce, hash, err
		}
	}

	return block.Nonce, hash, nil
}
