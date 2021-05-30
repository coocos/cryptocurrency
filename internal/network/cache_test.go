package network

import (
	"reflect"
	"testing"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

func TestCache(t *testing.T) {
	t.Run("Test reading block added to cache", func(t *testing.T) {
		cache := &BlockCache{}
		block := *blockchain.GenesisBlock()
		cache.AddBlock(block)
		readBlock := cache.ReadLastBlock()
		if !reflect.DeepEqual(block, readBlock) {
			t.Error("Cache returned wrong block")
		}
	})
}
