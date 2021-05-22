package keys

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"log"
	"os"
)

// KeyPair is a valid Ed25519 key pair
type KeyPair struct {
	PrivateKey      ed25519.PrivateKey
	PublicKey       ed25519.PublicKey
	PublicKeyBase64 string
}

// NewKeyPair returns a random new key pair
func NewKeyPair() *KeyPair {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v\n", err)
	}
	return &KeyPair{
		privateKey,
		publicKey,
		base64.StdEncoding.EncodeToString(publicKey),
	}
}

// LoadKeyPair loads key pair from private key file
func LoadKeyPair(privateKeyPath string) (*KeyPair, error) {
	privateKeySeed, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	privateKey := ed25519.NewKeyFromSeed(privateKeySeed)
	publicKey, ok := privateKey.Public().(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("Failed to convert public key")
	}
	return &KeyPair{
		privateKey,
		publicKey,
		base64.StdEncoding.EncodeToString(publicKey),
	}, nil
}
