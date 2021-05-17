package blockchain

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"
)

func TestBlock(t *testing.T) {
	t.Run("Test hashing a block", func(t *testing.T) {
		block := NewBlock(0, nil, nil)
		block.Time = time.Date(2021, time.January, 1, 6, 0, 0, 0, time.UTC)

		hash := block.ComputeHash()
		expectedHash, _ := hex.DecodeString("585d1618d011fcf9bb20ae920ebd29f69caa27f62a0635043e1d4d659696b883")

		if !bytes.Equal(hash, expectedHash) {
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

		block := NewBlock(chain.LastBlock().Number+1, chain.LastBlock().Hash, nil)
		for !block.IsValid() {
			block = NewBlock(chain.LastBlock().Number+1, chain.LastBlock().Hash, nil)
		}

		err := chain.AddBlock(block)
		if err != nil {
			t.Errorf("Failed to add block to blockhain: %s\n", err)
		}
	})
}
