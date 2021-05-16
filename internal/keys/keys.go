package keys

import (
	"crypto/ed25519"
	"crypto/x509"
	"errors"
	"log"
	"os"
)

// KeyPair is a valid Ed25519 key pair
type KeyPair struct {
	PrivateKey        ed25519.PrivateKey
	PublicKey         ed25519.PublicKey
	EncodedPrivateKey []byte
	EncodedPublicKey  []byte
}

// NewKeyPair returns a random new key pair
func NewKeyPair() *KeyPair {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("Failed to generate key pair: %v\n", err)
	}
	encodedPublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatalf("Failed to encode public key: %v\n", err)
	}
	encodedPrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Failed to encode private key: %v\n", err)
	}
	return &KeyPair{
		privateKey,
		publicKey,
		encodedPrivateKey,
		encodedPublicKey,
	}
}

func loadPrivateKey(privateKeyPath string) (ed25519.PrivateKey, error) {
	encodedPrivateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(encodedPrivateKey)
	if err != nil {
		return nil, err
	}
	edPrivateKey, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("Private key is not a valid Ed25519 private key")
	}
	return edPrivateKey, nil
}

func loadPublicKey(publicKeyPath string) (ed25519.PublicKey, error) {
	encodedPublicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := x509.ParsePKIXPublicKey(encodedPublicKey)
	if err != nil {
		return nil, err
	}
	edPublicKey, ok := publicKey.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("Public key is not a valid Ed25519 public key")
	}
	return edPublicKey, nil
}

// LoadKeyPair loads key pair from disk
func LoadKeyPair(privateKeyPath string, publicKeyPath string) (*KeyPair, error) {
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		log.Printf("Failed to read private key %v: %v\n", privateKeyPath, err)
		return nil, err
	}
	publicKey, err := loadPublicKey(publicKeyPath)
	if err != nil {
		return nil, err
	}
	encodedPublicKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatalf("Failed to encode public key: %v\n", err)
	}
	encodedPrivateKey, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Fatalf("Failed to encode public key: %v\n", err)
	}
	return &KeyPair{
		privateKey,
		publicKey,
		encodedPrivateKey,
		encodedPublicKey,
	}, nil
}
