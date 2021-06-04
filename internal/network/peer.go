package network

import (
	"log"
	"sync"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

// Peers contains all the peers a node knows about
type Peers struct {
	sync.RWMutex
	hosts map[string]bool
}

// Add adds a new peer node and shakes hands with it
func (p *Peers) Add(address string) {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.hosts[address]; ok {
		return
	}

	client := NodeClient{address}
	if err := client.Greet(); err != nil {
		log.Printf("Peer failed to respond, dropping it: %v\n", err)
		return
	}
	if p.hosts == nil {
		p.hosts = make(map[string]bool)
	}
	p.hosts[address] = true
}

// GetBlocks gets all blocks from peer node
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

// BroadcastBlock sends block to all known peer ndoes
func (p *Peers) BroadcastBlock(block blockchain.Block) {
	p.RLock()
	defer p.RUnlock()

	done := make(chan bool)
	for host := range p.hosts {
		go func(host string, block blockchain.Block, done chan<- bool) {
			client := NodeClient{host}
			err := client.SendBlock(block)
			if err != nil {
				log.Printf("Failed to broadcast block to %s: %v\n", host, err)
			}
			done <- true
		}(host, block, done)
	}
	for range p.hosts {
		<-done
	}
}
