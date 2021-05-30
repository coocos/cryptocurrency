package network

import (
	"sync"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

// BlockCache is a synchronized cache for blocks
type BlockCache struct {
	sync.RWMutex
	blocks []blockchain.Block
}

// AddBlock adds a block to the cache
func (b *BlockCache) AddBlock(block blockchain.Block) {
	b.Lock()
	b.blocks = append(b.blocks, block)
	b.Unlock()
}

// ReadBlock returns a block from the cache
func (b *BlockCache) ReadLastBlock() blockchain.Block {
	b.RLock()
	block := b.blocks[len(b.blocks)-1]
	b.RUnlock()
	return block
}

// ReadBlocks returns a channel for iterating over all the blocks in the cache
func (b *BlockCache) ReadBlocks() <-chan blockchain.Block {
	blocks := make(chan blockchain.Block)
	go func() {
		b.RLock()
		defer b.RUnlock()
		for _, block := range b.blocks {
			blocks <- block
		}
		close(blocks)
	}()
	return blocks
}
