package blockchain

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"
)

func TestBlock(t *testing.T) {
	t.Run("Test hashing a block", func(t *testing.T) {
		genesisBlock := GenesisBlock()
		block := NewBlock(genesisBlock.Number+1, genesisBlock.Hash, nil)
		block.Time = time.Date(2021, time.January, 1, 6, 0, 0, 0, time.UTC)
		block.Nonce = 1

		hash := block.ComputeHash()
		expectedHash, _ := hex.DecodeString("0323204fb3d384837511fef3da59fc5324655e0376a2d3e624c11d418b873a45")

		if !bytes.Equal(hash, expectedHash) {
			t.Errorf("Block hash %x differs from expected %x\n", hash, expectedHash)
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
