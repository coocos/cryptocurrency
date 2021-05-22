package blockchain

import (
	"testing"

	"github.com/coocos/cryptocurrency/internal/keys"
)

func TestTransaction(t *testing.T) {
	t.Run("Test valid signature", func(t *testing.T) {
		senderKeyPair := keys.NewKeyPair()
		receiverKeyPair := keys.NewKeyPair()

		transaction := Transaction{
			Sender:   senderKeyPair.PublicKey,
			Receiver: receiverKeyPair.PublicKey,
			Amount:   10,
		}
		transaction.Sign(senderKeyPair.PrivateKey)

		if !transaction.ValidSignature() {
			t.Errorf("Expected signature to be valid")
		}
	})
	t.Run("Test invalid signature", func(t *testing.T) {
		senderKeyPair := keys.NewKeyPair()
		receiverKeyPair := keys.NewKeyPair()

		transaction := Transaction{
			Sender:   senderKeyPair.PublicKey,
			Receiver: receiverKeyPair.PublicKey,
			Amount:   10,
		}
		transaction.Sign(receiverKeyPair.PrivateKey)

		if transaction.ValidSignature() {
			t.Errorf("Expected signature to be invalid")
		}
	})
	t.Run("Test identifying coinbase transactions", func(t *testing.T) {
		minerKeyPair := keys.NewKeyPair()

		transaction := Transaction{
			Sender:   nil,
			Receiver: minerKeyPair.PublicKey,
			Amount:   10,
		}
		if !transaction.IsCoinbase() {
			t.Error("Coinbase transaction not identified as coinbase transaction")
		}
		if !transaction.ValidSignature() {
			t.Error("Coinbase transactions signatures are always considered valid")
		}
	})
}
