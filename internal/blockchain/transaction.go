package blockchain

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	CoinbaseTransactionAmount = 10
)

// Transaction represents an individual transaction
type Transaction struct {
	Sender    []byte `json:"sender"`
	Receiver  []byte `json:"receiver"`
	Amount    uint   `json:"amount"`
	Nonce     uint   `json:"nonce"`
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
func NewTransaction(sender ed25519.PublicKey, receiver ed25519.PublicKey, amount uint, nonce uint) *Transaction {
	return &Transaction{
		sender,
		receiver,
		amount,
		nonce,
		nil,
	}
}

// Bytes returns the transaction as bytes
func (t *Transaction) Bytes() ([]byte, error) {
	// Omit the signature since the signature is used to sign this
	copy := Transaction{
		Sender:   t.Sender,
		Receiver: t.Receiver,
		Amount:   t.Amount,
		Nonce:    t.Nonce,
	}

	bytes, err := json.Marshal(copy)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Sign signs the transaction using the given key and returns the signature
func (t *Transaction) Sign(privateKey ed25519.PrivateKey) ([]byte, error) {
	bytes, err := t.Bytes()
	if err != nil {
		return nil, err
	}
	signature := ed25519.Sign(privateKey, bytes)
	t.Signature = signature
	return signature, nil
}

// IsCoinBase tells whether the transaction is a coinbase transaction
func (t *Transaction) IsCoinbase() bool {
	return t.Sender == nil && t.Receiver != nil && t.Amount == CoinbaseTransactionAmount
}

// ValidSignature indicates whether the transaction signature is valid
func (t *Transaction) ValidSignature() bool {
	// Coinbase transactions do not have signatures
	if t.Sender == nil {
		return t.IsCoinbase()
	}

	bytes, err := t.Bytes()
	if err != nil {
		return false
	}
	return ed25519.Verify(t.Sender, bytes, t.Signature)
}

// CoinbaseTransaction contructs a coinbase transaction
func CoinbaseTransactionTo(receiver ed25519.PublicKey) Transaction {
	return Transaction{
		Sender:   nil,
		Receiver: receiver,
		Amount:   CoinbaseTransactionAmount,
	}
}
