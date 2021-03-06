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
		block := NewBlock(genesisBlock.Number+1, genesisBlock.Hash, nil, 1)
		block.Time = time.Date(2021, time.January, 1, 6, 0, 0, 0, time.UTC)

		hash := block.ComputeHash()
		expectedHash, _ := hex.DecodeString("0323204fb3d384837511fef3da59fc5324655e0376a2d3e624c11d418b873a45")

		if !bytes.Equal(hash, expectedHash) {
			t.Errorf("Block hash %x differs from expected %x\n", hash, expectedHash)
		}
	})
}
