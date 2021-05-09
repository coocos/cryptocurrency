package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"
)

// Blockchain represents a full blockchain
type Blockchain struct {
	chain []*Block
}

// NewBlockchain returns a new blockchain with a genesis block
func NewBlockchain() *Blockchain {
	blockchain := Blockchain{}

	// FIXME: Genesis block is not valid
	genesisBlock := NewBlock(0, []byte{})
	blockchain.AddBlock(genesisBlock)

	return &blockchain
}

// LastBlock returns the last block in the blockchain
func (b *Blockchain) LastBlock() *Block {
	if len(b.chain) > 0 {
		return b.chain[len(b.chain)-1]
	}
	return nil
}

// AddBlock adds block to blockchain or returns an error if the block is not valid
func (b *Blockchain) AddBlock(block *Block) error {
	previous := b.LastBlock()
	if previous == nil {
		log.Printf("Adding genesis block: %+v\n", block)
		b.chain = append(b.chain, block)
		return nil
	}
	if !block.IsValid() {
		return errors.New("New block not valid according to proof-of-work")
	}
	if block.Number != previous.Number+1 || !bytes.Equal(block.PreviousHash, previous.Hash) {
		return errors.New("New block does not follow the last block in blockchain")
	}
	b.chain = append(b.chain, block)
	return nil
}

// Mine executes the proof-of-work algorithm to mine for a valid block
func (b *Blockchain) Mine() {
	for {
		block := NewBlock(b.LastBlock().Number+1, b.LastBlock().PreviousHash)
		if !block.IsValid() {
			continue
		}
		log.Printf("Found valid block: %+v\n", block)
		b.AddBlock(block)
		break
	}
}

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

// Transaction represents a transaction within a block
type Transaction struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   int    `json:"amount"`
}
