package network

import (
	"log"
	"sync"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

// Peers contains all the peers a node knows about
type Peers struct {
	sync.RWMutex
	peers map[string]bool
}

// Add adds a new peer and shakes hands with it
func (p *Peers) Add(address string) {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.peers[address]; ok {
		return
	}

	client := NodeClient{address}
	if err := client.Greet(); err != nil {
		log.Printf("Peer failed to respond, dropping it: %v\n", err)
		return
	}
	if p.peers == nil {
		p.peers = make(map[string]bool)
	}
	p.peers[address] = true
}

// GetBlocks gets all blocks from peer
func (p *Peers) GetBlocks(address string) ([]blockchain.Block, error) {
	p.RLock()
	defer p.RUnlock()
	client := NodeClient{address}
	blocks, err := client.GetBlocks()
	if err != nil {
		return []blockchain.Block{}, err
	}
	return blocks, nil
}
