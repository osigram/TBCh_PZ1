package storage

import "PZ1/internal/domain"

type Storage interface {
	Blockchain() domain.Blockchain
	AddBlock(block *domain.Block) error
	SetBlockchain(blockchain domain.Blockchain)
}
