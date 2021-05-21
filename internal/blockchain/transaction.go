package blockchain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Transaction represents an individual transaction
type Transaction struct {
	Sender    []byte `json:"sender"`
	Receiver  []byte `json:"receiver"`
	Amount    uint   `json:"amount"`
	Signature []byte `json:"signature"`
}

// String returns the string representation of a transaction
func (t Transaction) String() string {
	sender := base64.StdEncoding.EncodeToString(t.Sender)
	receiver := base64.StdEncoding.EncodeToString(t.Receiver)
	if t.Sender == nil {
		return fmt.Sprintf("Transaction: %d coins to miner %s", t.Amount, receiver)
	}
	return fmt.Sprintf("Transaction: %d coins from %s to %s", t.Amount, sender, receiver)
}

// NewTransaction returns a new unsigned transaction
func NewTransaction(sender []byte, receiver []byte, amount uint) *Transaction {
	// FIXME: Validate that both sender and receiver are proper public keys
	return &Transaction{
		sender,
		receiver,
		amount,
		nil,
	}
}

// ComputeHash sets the hash for transaction and returns it
func (t *Transaction) ComputeHash() ([]byte, error) {
	// Hash should not include the signature
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

	return hash, nil
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

// IsCoinBase tells whether the transaction is a coinbase transaction
func (t *Transaction) IsCoinbase() bool {
	return t.Sender == nil && t.Receiver != nil && t.Amount == 10
}

// ValidSignature indicates whether the transaction signature is valid
func (t *Transaction) ValidSignature() bool {
	// Coinbase transactions do not have signatures
	if t.Sender == nil {
		return t.IsCoinbase()
	}

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
