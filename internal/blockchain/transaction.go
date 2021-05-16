package blockchain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
)

// Transaction represents an individual transaction
type Transaction struct {
	Hash      []byte `json:"hash"`
	Sender    []byte `json:"sender"`
	Receiver  []byte `json:"receiver"`
	Amount    int    `json:"amount"`
	Signature []byte `json:"signature"`
}

// ComputeHash sets the hash for transaction and returns it
func (t *Transaction) ComputeHash() ([]byte, error) {
	// Hash should not include the hash itself or the signature
	copy := Transaction{
		Sender:   t.Sender,
		Receiver: t.Receiver,
		Amount:   t.Amount,
	}

	bytes, err := json.Marshal(copy)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write(bytes)
	hash := h.Sum(nil)

	t.Hash = hash
	return t.Hash, nil
}

// Sign signs the transaction using the given key and returns the signature
func (t *Transaction) Sign(privateKey ed25519.PrivateKey) ([]byte, error) {
	hash, err := t.ComputeHash()
	if err != nil {
		return nil, err
	}
	signature := ed25519.Sign(privateKey, hash)
	t.Signature = signature
	return signature, nil
}

// ValidSignature indicates whether the transaction signature is valid
func (t *Transaction) ValidSignature() bool {
	parsedKey, err := x509.ParsePKIXPublicKey(t.Sender)
	if err != nil {
		return false
	}
	key, ok := parsedKey.(ed25519.PublicKey)
	if !ok {
		return false
	}
	transactionHash, err := t.ComputeHash()
	if err != nil {
		return false
	}
	return ed25519.Verify(key, transactionHash, t.Signature)
}
