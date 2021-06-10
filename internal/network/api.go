package network

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coocos/cryptocurrency/internal/blockchain"
	"github.com/coocos/cryptocurrency/internal/config"
)

// Api runs the HTTP API for interacting with the node
type Api struct {
	cache  *BlockCache
	events chan<- interface{}
}

// NewApi returns a new instance of the API server
func NewApi(events chan<- interface{}) *Api {
	return &Api{
		&BlockCache{},
		events,
	}
}

func (a *Api) updateCache(block blockchain.Block) {
	a.cache.AddBlock(block)
}

// Serve starts the API
func (a *Api) Serve() error {
	// Returns blocks from the blockchain
	http.HandleFunc("/api/v1/blockchain/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		blocks := []blockchain.Block{}
		for block := range a.cache.ReadBlocks() {
			blocks = append(blocks, block)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blocks)
	})
	// Receives new blocks from other nodes
	http.HandleFunc("/api/v1/block/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var block NewBlock
		if err := json.NewDecoder(r.Body).Decode(&block); err != nil {
			http.Error(w, "Request is not valid JSON", http.StatusBadRequest)
			return
		}
		a.events <- block
		w.WriteHeader(http.StatusAccepted)
	})
	http.HandleFunc("/api/v1/accounts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		blocks := []*blockchain.Block{}
		for block := range a.cache.ReadBlocks() {
			blocks = append(blocks, &block)
		}
		accounts := blockchain.AccountsFromBlockchain(blocks).ListAccounts()
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(accounts)
		if err != nil {
			log.Println("Failed to serialize accounts", err)
		}
	})
	// Receives notifications of new peer nodes
	http.HandleFunc("/api/v1/peer/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var peer NewPeer
		if err := json.NewDecoder(r.Body).Decode(&peer); err != nil {
			http.Error(w, "Request is not valid JSON", http.StatusBadRequest)
			return
		}
		a.events <- peer
		w.Write(nil)
	})
	bindHost := config.BindHost()
	log.Println("Listening for API requests at", bindHost)
	return http.ListenAndServe(bindHost, nil)
}
