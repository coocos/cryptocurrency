package network

import (
	"encoding/json"
	"net/http"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

// Api runs the HTTP API for interacting with the node
type Api struct {
	cache             *BlockCache
	unconfirmedBlocks chan<- blockchain.Block
}

// NewApi returns a new instance of the API server
func NewApi(cache *BlockCache, unconfirmedBlocks chan<- blockchain.Block) *Api {
	return &Api{
		cache,
		unconfirmedBlocks,
	}
}

// Serve starts the API
func (a *Api) Serve() {
	// Returns blocks from the blockchain
	http.HandleFunc("/api/v1/blockchain/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		block := a.cache.ReadLastBlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(block)
	})
	// Receives new blocks from other nodes
	http.HandleFunc("/api/v1/block/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var block blockchain.Block
		err := json.NewDecoder(r.Body).Decode(&block)
		if err != nil {
			http.Error(w, "Block is not valid JSON", http.StatusBadRequest)
			return
		}
		a.unconfirmedBlocks <- block
		w.WriteHeader(http.StatusAccepted)
	})
	http.ListenAndServe("localhost:8080", nil)
}