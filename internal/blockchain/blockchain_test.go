package blockchain

import (
	"testing"
	"time"
)

func TestBlockChain(t *testing.T) {
	t.Run("Test that blockchain includes genesis block", func(t *testing.T) {
		chain := NewBlockchain()

		if chain.LastBlock() == nil {
			t.Errorf("Blockchain has no genesis block\n")
		}
	})
	t.Run("Test adding a new valid block", func(t *testing.T) {
		chain := NewBlockchain()

		block := NewBlock(chain.LastBlock().Number+1, chain.LastBlock().Hash, nil, 14859)
		block.Time = time.Date(2021, time.May, 1, 6, 0, 0, 0, time.UTC)
		block.Hash = block.ComputeHash()

		err := chain.AddBlock(block)
		if err != nil {
			t.Errorf("Failed to add block to blockhain: %s\n", err)
		}
	})
}
