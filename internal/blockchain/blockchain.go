package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"time"
)

// Blockchain represents a full blockchain
type Blockchain struct {
	chain []*Block
}

func (b *Blockchain) LastBlock() *Block {
	if len(b.chain) > 0 {
		return b.chain[len(b.chain)-1]
	}
	return nil
}

func (b *Blockchain) AddBlock(block *Block) {
	b.chain = append(b.chain, block)
}

// Block is an individual block in the blockchain
type Block struct {
	Number       int            `json:"number"`
	Time         time.Time      `json:"time"`
	Transactions []*Transaction `json:"transactions"`
	PreviousHash []byte         `json:"previousHash"`
	Hash         []byte         `json:"hash"`
}

// NewBlock creates a new block
func NewBlock(number int, previousHash []byte) *Block {
	return &Block{
		Number:       number,
		Time:         time.Now().UTC(),
		PreviousHash: previousHash,
	}
}

func (b *Block) addTransaction(transaction *Transaction) {
	b.Transactions = append(b.Transactions, transaction)
}

// Compute hash computes the hash for the block
func (b *Block) ComputeHash() []byte {

	// Hashing the block needs to exclude the hash itself
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

// Transactions represents a transaction between sender and receiver
type Transaction struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   int    `json:"amount"`
}
