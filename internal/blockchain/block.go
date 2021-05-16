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
	PreviousHash string         `json:"previousHash"`
	Hash         string         `json:"hash"`
}

// NewBlock creates a new block
func NewBlock(number int, previousHash string) *Block {
	block := Block{
		Number:       number,
		Time:         time.Now().UTC(),
		PreviousHash: previousHash,
		Nonce:        rand.Int63(),
	}
	block.Hash = block.ComputeHash()
	return &block
}

// GenesisBlock returns the fixed first block in the blockchain
func GenesisBlock() *Block {
	block := Block{
		Number:       0,
		Time:         time.Date(2021, time.May, 1, 6, 0, 0, 0, time.UTC),
		PreviousHash: "",
		Nonce:        4923246119299551551,
		Hash:         "0000b049046735988b782ddd65ee6f49ec4d5501e84bb229b53d52dded20f5c0",
	}
	return &block
}

// AddTransaction adds a transaction to block
func (b *Block) AddTransaction(transaction *Transaction) {
	b.Transactions = append(b.Transactions, transaction)
}

// ComputeHash computes the hash for the block
func (b *Block) ComputeHash() string {
	// Exclude the hash field itself when hashing the block
	copy := Block{
		Number:       b.Number,
		Time:         b.Time,
		Transactions: b.Transactions,
		Nonce:        b.Nonce,
		PreviousHash: b.PreviousHash,
	}
	bytes, err := json.Marshal(copy)
	if err != nil {
		log.Fatalf("Failed to convert block to bytes: %s\n", err)
	}

	hash := sha256.New()
	hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

// IsValid indicates if the block hash is valid
func (b *Block) IsValid() bool {
	return b.Hash[:4] == "0000"
}
