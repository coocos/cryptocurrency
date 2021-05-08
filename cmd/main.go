package main

import (
	"encoding/json"
	"log"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

func main() {

	block := blockchain.NewBlock(0, []byte{})
	serialized, err := json.Marshal(block)
	if err != nil {
		log.Fatalf("Failed to serialize block to JSON: %s\n", err)
	}
	log.Println(string(serialized))

}
