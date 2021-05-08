package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log"
	"time"
)

// Blockchain represents a full blockchain
type Blockchain struct {
	chain []*Block
}

// NewBlockchain returns a new blockchain with a genesis block
func NewBlockchain() *Blockchain {
	blockchain := Blockchain{}

	genesisBlock := NewBlock(0, []byte{})
	blockchain.AddBlock(genesisBlock)

	return &blockchain
}

func (b *Blockchain) LastBlock() *Block {
	if len(b.chain) > 0 {
		return b.chain[len(b.chain)-1]
	}
	return nil
}

func (b *Blockchain) AddBlock(block *Block) error {
	previous := b.LastBlock()
	if previous == nil {
		log.Printf("Adding genesis block: %+v\n", block)
		b.chain = append(b.chain, block)
		return nil
	}
	if block.Number != previous.Number+1 || !bytes.Equal(block.PreviousHash, previous.Hash) {
		return errors.New("New block does not follow the latest block in blockchain")
	}
	b.chain = append(b.chain, block)
	return nil
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
	block := Block{
		Number:       number,
		Time:         time.Now().UTC(),
		PreviousHash: previousHash,
	}
	block.Hash = block.ComputeHash()
	return &block
}

func (b *Block) addTransaction(transaction *Transaction) {
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

// Transaction represents a transaction within a block
type Transaction struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   int    `json:"amount"`
}
