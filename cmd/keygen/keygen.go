// Tool for generating an Ed25519 key pair
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/coocos/cryptocurrency/internal/keys"
)

type Options struct {
	privateKeyFile string
	publicKeyFile  string
}

func parseFlags() Options {
	options := Options{}
	flag.StringVar(&options.privateKeyFile, "private", "private.key", "Private key file name")
	flag.StringVar(&options.publicKeyFile, "public", "public.key", "Public key file name")
	flag.Parse()
	return options
}

func writeKeyPairToFile(keyPair keys.KeyPair, options Options) error {
	if err := os.WriteFile(options.privateKeyFile, keyPair.PrivateKey.Seed(), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(options.publicKeyFile, keyPair.PublicKey, 0644); err != nil {
		return err
	}
	return nil
}

func main() {
	options := parseFlags()

	fmt.Printf("⏳ Generating key pair %s and %s...\n", options.privateKeyFile, options.publicKeyFile)
	keyPair := keys.NewKeyPair()
	if err := writeKeyPairToFile(*keyPair, options); err != nil {
		log.Fatalf("Failed to write key pair to file: %v\n", err)
	}
	fmt.Println("✨ Done!")
}
