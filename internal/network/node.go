package network

import (
	"log"
	"os"

	"github.com/coocos/cryptocurrency/internal/blockchain"
	"github.com/coocos/cryptocurrency/internal/keys"
)

// Node represents the node running the blockchain
type Node struct {
	chain *blockchain.Blockchain
	api   *Api
	peers *Peers
}

// NewNode returns a new node which mines blocks using the given key pair
func NewNode(keyPair *keys.KeyPair) *Node {
	chain := blockchain.NewBlockchain(keyPair)
	peers := &Peers{}
	events := eventBus(chain, peers)
	api := NewApi(events)
	return &Node{
		chain: chain,
		api:   api,
		peers: peers,
	}
}

// Start starts the node
func (n *Node) Start() {
	n.serve()
	if seedHost, ok := os.LookupEnv("NODE_SEED_HOST"); ok {
		log.Println("Syncing blockchain via", seedHost)
		n.peers.Add(seedHost)
		blocks, err := n.peers.GetBlocks(seedHost)
		if err != nil {
			log.Println("Failed to sync blockchain using seed node:", err)
		} else {
			for _, block := range blocks {
				n.chain.SubmitExternalBlock(&block)
			}
		}
	}
	n.mine()
}

func eventBus(chain *blockchain.Blockchain, peers *Peers) chan<- interface{} {
	events := make(chan interface{})
	go func() {
		for event := range events {
			switch e := event.(type) {
			case NewBlock:
				chain.SubmitExternalBlock(&e.Block)
			case NewPeer:
				log.Println("Node @", e.Address, "sent greeting")
				peers.Add(e.Address)
			default:
				log.Fatalf("Received an unknown event: %v\n", event)

			}
		}
	}()
	return events
}

func (n *Node) serve() {
	go func() {
		if err := n.api.Serve(); err != nil {
			log.Fatalln("API serving failed:", err)
		}
	}()
}

func (n *Node) mine() {
	for {
		block := n.chain.MineBlock()
		n.api.updateCache(block)
		n.peers.BroadcastBlock(block)
	}
}
