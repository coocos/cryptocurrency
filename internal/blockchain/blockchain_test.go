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

		hash = block.ComputeHash()
		expectedHash = "dabe48d04485375a2f43608c57586670ba1d708568aad459a21bd7615283f9df"

		if hex.EncodeToString(hash) != expectedHash {
			t.Errorf("Block hash %x differs from expected %s\n", hash, expectedHash)
		}

	})

}
