package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Block is an individual block in the blockchain
type Block struct {
	Number       int           `json:"number"`
	Time         time.Time     `json:"time"`
	Transactions []Transaction `json:"transactions"`
	Nonce        int           `json:"nonce"`
	PreviousHash []byte        `json:"previousHash"`
	Hash         []byte        `json:"hash"`
}

// String returns the string representation of a block
func (b Block) String() string {
	return fmt.Sprintf("Block %d %x transactions: %d", b.Number, b.Hash, len(b.Transactions))
}

// NewBlock creates a new block
func NewBlock(number int, previousHash []byte, transactions []Transaction, nonce int) *Block {
	block := Block{
		Number:       number,
		Time:         time.Now().UTC(),
		Transactions: transactions,
		PreviousHash: previousHash,
		Nonce:        nonce,
	}
	block.Hash = block.ComputeHash()
	return &block
}

// GenesisBlock returns the fixed first block in the blockchain
func GenesisBlock() *Block {
	block := Block{
		Number:       0,
		Time:         time.Date(2021, time.May, 1, 6, 0, 0, 0, time.UTC),
		PreviousHash: nil,
		Nonce:        3999606801082803789,
	}
	hash, _ := hex.DecodeString("000002be9afbfdaa977028a51d10bd590f9b56b03c3f570b8723e3809dc439ba")
	block.Hash = hash
	if !block.IsValid() {
		log.Fatalln("Genesis block is not valid")
	}
	return &block
}

// ComputeHash computes the hash for the block
func (b *Block) ComputeHash() []byte {
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
		log.Fatalf("Failed to hash block: %v\n", err)
	}

	hash := sha256.New()
	hash.Write(bytes)
	return hash.Sum(nil)
}

// IsValid indicates if the block hash is valid
func (b *Block) IsValid() bool {
	return bytes.Equal(b.ComputeHash(), b.Hash) &&
		hex.EncodeToString(b.Hash)[:4] == "0000"
}
