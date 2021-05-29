package network

import (
	"encoding/json"
	"net/http"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

// Api runs the HTTP API for interacting with the node
type Api struct {
	blocks BlockSource
}

// NewApi returns a new instance of the API server
func NewApi(blocks BlockSource) *Api {
	return &Api{
		blocks,
	}
}

// Serve starts the API
func (a *Api) Serve() {
	http.HandleFunc("/api/v1/blockchain/", func(w http.ResponseWriter, r *http.Request) {
		block := a.blocks.ReadBlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(block)
	})
	http.HandleFunc("/api/v1/block/", func(w http.ResponseWriter, r *http.Request) {
		var block blockchain.Block
		err := json.NewDecoder(r.Body).Decode(&block)
		if err != nil {
			http.Error(w, "Block is not valid JSON", http.StatusBadRequest)
			return
		}
		a.blocks.SubmitBlock(block)
		w.WriteHeader(http.StatusAccepted)
	})
	http.ListenAndServe("localhost:8080", nil)
}
