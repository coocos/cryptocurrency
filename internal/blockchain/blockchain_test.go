package blockchain

import (
	"testing"
	"time"
)

func TestBlock(t *testing.T) {
	t.Run("Test hashing a block", func(t *testing.T) {
		block := NewBlock(0, "")
		block.Time = time.Date(2021, time.January, 1, 6, 0, 0, 0, time.UTC)

		hash := block.ComputeHash()
		expectedHash := "5bf97525874b93320a47e21030dbc117696121f0fd36a9afdfcdd0fddf817e26"

		if hash != expectedHash {
			t.Errorf("Block hash %s differs from expected %s\n", hash, expectedHash)
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
	t.Run("Test adding a valid block", func(t *testing.T) {
		chain := NewBlockchain()

		block := NewBlock(chain.LastBlock().Number+1, chain.LastBlock().Hash)
		for !block.IsValid() {
			block = NewBlock(chain.LastBlock().Number+1, chain.LastBlock().Hash)
		}

		err := chain.AddBlock(block)
		if err != nil {
			t.Errorf("Failed to add block to blockhain: %s\n", err)
		}
	})
}
