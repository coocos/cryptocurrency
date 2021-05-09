package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

// Block is an individual block in the blockchain
type Block struct {
	Number       int            `json:"number"`
	Time         time.Time      `json:"time"`
	Transactions []*Transaction `json:"transactions"`
	Nonce        int64          `json:"nonce"`
	PreviousHash []byte         `json:"previousHash"`
	Hash         []byte         `json:"hash"`
}

// NewBlock creates a new block
func NewBlock(number int, previousHash []byte) *Block {
	block := Block{
		Number:       number,
		Time:         time.Now().UTC(),
		PreviousHash: previousHash,
		Nonce:        rand.Int63(),
	}
	block.Hash = block.ComputeHash()
	return &block
}

// AddTransaction adds a transaction to block
func (b *Block) AddTransaction(transaction *Transaction) {
	b.Transactions = append(b.Transactions, transaction)
}

// ComputeHash computes the hash for the block
func (b *Block) ComputeHash() []byte {
	// Exclude the hash field itself when hashing the block
	copy := b
	copy.Hash = []byte{}
	bytes, err := json.Marshal(copy)
	if err != nil {
		log.Fatalf("Failed to convert block to bytes: %s\n", err)
	}

	hash := sha256.New()
	hash.Write(bytes)
	return hash.Sum(nil)
}

// IsValid indicates if the block hash is valid
func (b *Block) IsValid() bool {
	return hex.EncodeToString(b.Hash)[:4] == "0000"
}
