package blockchain

import (
	"encoding/hex"
	"testing"
	"time"
)

func TestBlock(t *testing.T) {
	t.Run("Test hashing a block", func(t *testing.T) {
		block := NewBlock(0, []byte{})
		block.Time = time.Date(2021, time.January, 1, 6, 0, 0, 0, time.UTC)

		hash := block.ComputeHash()
		expectedHash := "dabe48d04485375a2f43608c57586670ba1d708568aad459a21bd7615283f9df"

		if hex.EncodeToString(hash) != expectedHash {
			t.Errorf("Block hash %x differs from expected %s\n", hash, expectedHash)
		}
	})
}

func TestBlockChain(t *testing.T) {
	t.Run("Test creating a blockchain", func(t *testing.T) {
		chain := NewBlockchain()

		if chain.LastBlock() == nil {
			t.Errorf("Blockchain has no genesis block\n")
		}
	})
	t.Run("Test adding a block", func(t *testing.T) {
		chain := NewBlockchain()

		block := NewBlock(chain.LastBlock().Number+1, chain.LastBlock().Hash)
		err := chain.AddBlock(block)
		if err != nil {
			t.Errorf("Failed to add block to blockhain: %s\n", err)
		}
	})
}
