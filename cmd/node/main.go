package main

import (
	"log"
	"os"

	"github.com/coocos/cryptocurrency/internal/blockchain"
	"github.com/coocos/cryptocurrency/internal/keys"
)

func minerKeyPair() *keys.KeyPair {
	privateKeyPath := os.Getenv("NODE_PRIVATE_KEY")
	publicKeyPath := os.Getenv("NODE_PUBLIC_KEY")
	keyPair, err := keys.LoadKeyPair(privateKeyPath, publicKeyPath)
	if err != nil {
		log.Fatalf("Failed to load key pair, unable to sign transactions: %v\n", err)

	}
	return keyPair
}

func main() {

	chain := blockchain.NewBlockchain(minerKeyPair())
	for {
		chain.MineBlock()
	}

}
