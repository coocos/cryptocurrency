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
			Sender:   senderKeyPair.EncodedPublicKey,
			Receiver: receiverKeyPair.EncodedPublicKey,
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
			Sender:   senderKeyPair.EncodedPublicKey,
			Receiver: receiverKeyPair.EncodedPublicKey,
			Amount:   10,
		}
		transaction.Sign(receiverKeyPair.PrivateKey)

		if transaction.ValidSignature() {
			t.Errorf("Expected signature to be invalid")
		}
	})
}
