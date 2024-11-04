package inmemory

import (
	"PZ1/internal/domain"
	"crypto/ecdsa"
	"sync"
)

type BlockchainStorage struct {
	blockchain *domain.Blockchain
	mx         sync.RWMutex
}

func MustNewBlockchainStorage(key *ecdsa.PrivateKey) *BlockchainStorage {
	firstTransaction, err := domain.NewRegisterTransaction(key)
	if err != nil {
		panic(err)
	}

	firstBlock, err := domain.NewBlock(nil, *firstTransaction)
	if err != nil {
		panic(err)
	}

	blockchain := domain.Blockchain([]domain.Block{*firstBlock})

	return &BlockchainStorage{blockchain: &blockchain}
}

func (s *BlockchainStorage) Blockchain() domain.Blockchain {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return *s.blockchain
}

func (s *BlockchainStorage) AddBlock(block *domain.Block) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if err := block.Validate(*s.blockchain); err != nil {
		return err
	}

	return s.blockchain.Add(*block)
}

func (s *BlockchainStorage) SetBlockchain(blockchain domain.Blockchain) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.blockchain = &blockchain
}
