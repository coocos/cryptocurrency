package blockchain

import (
	"crypto/ed25519"
	"crypto/x509"
	"testing"
)

type KeyPair struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	address    []byte
}

func generateKeyPair() KeyPair {
	publicKey, privateKey, _ := ed25519.GenerateKey(nil)
	bytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	return KeyPair{
		privateKey,
		publicKey,
		bytes,
	}
}

func TestTransaction(t *testing.T) {
	t.Run("Test valid signature", func(t *testing.T) {
		senderKeyPair := generateKeyPair()
		receiverKeyPair := generateKeyPair()

		transaction := Transaction{
			Sender:   senderKeyPair.address,
			Receiver: receiverKeyPair.address,
			Amount:   10,
		}
		transaction.Sign(senderKeyPair.privateKey)

		if !transaction.ValidSignature() {
			t.Errorf("Expected signature to be valid")
		}
	})
	t.Run("Test invalid signature", func(t *testing.T) {
		senderKeyPair := generateKeyPair()
		receiverKeyPair := generateKeyPair()

		transaction := Transaction{
			Sender:   senderKeyPair.address,
			Receiver: receiverKeyPair.address,
			Amount:   10,
		}
		transaction.Sign(receiverKeyPair.privateKey)

		if transaction.ValidSignature() {
			t.Errorf("Expected signature to be invalid")
		}
	})
}
