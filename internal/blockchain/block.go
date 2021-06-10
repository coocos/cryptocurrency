package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/bits"
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

const (
	maxTransactionsPerBlock = 64
	baseDifficulty          = 20
)

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
	genesisHash, _ := hex.DecodeString("000002be9afbfdaa977028a51d10bd590f9b56b03c3f570b8723e3809dc439ba")
	return &Block{
		Number:       0,
		Time:         time.Date(2021, time.May, 1, 6, 0, 0, 0, time.UTC),
		PreviousHash: nil,
		Nonce:        3999606801082803789,
		Hash:         genesisHash,
	}
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

// IsValid indicates if the block is valid
func (b *Block) IsValid(previous *Block) bool {
	if b.Number != previous.Number+1 || !bytes.Equal(b.PreviousHash, previous.Hash) {
		return false
	}
	if len(b.Transactions) < 1 || len(b.Transactions) > maxTransactionsPerBlock {
		return false
	}
	if !b.Transactions[0].IsCoinbase() {
		return false
	}
	if !bytes.Equal(b.Hash, b.ComputeHash()) {
		return false
	}
	return bits.LeadingZeros64(binary.BigEndian.Uint64(b.Hash)) > baseDifficulty
}
