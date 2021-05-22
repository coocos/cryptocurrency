package main

import (
	"flag"
	"log"

	"github.com/coocos/cryptocurrency/internal/blockchain"
	"github.com/coocos/cryptocurrency/internal/keys"
)

// Options passed as CLI flags
type Options struct {
	privateKey string
	publicKey  string
}

func parseArgs() Options {
	options := Options{}
	flag.StringVar(&options.privateKey, "private", "private.key", "private key path")
	flag.StringVar(&options.publicKey, "public", "public.key", "public key path")
	flag.Parse()
	return options
}

func loadMinerKeyPair(options Options) *keys.KeyPair {
	keyPair, err := keys.LoadKeyPair(options.privateKey)
	if err != nil {
		log.Fatalf("Failed to load key pair, unable to sign transactions: %v\n", err)

	}
	return keyPair
}

func main() {
	options := parseArgs()

	chain := blockchain.NewBlockchain(loadMinerKeyPair(options))
	for {
		chain.MineBlock()
	}
}
