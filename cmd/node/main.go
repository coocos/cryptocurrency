package main

import (
	"crypto/ed25519"
	"crypto/x509"
	"log"

	"github.com/coocos/cryptocurrency/internal/blockchain"
)

func main() {

	senderPub, senderPriv, _ := ed25519.GenerateKey(nil)
	senderAddress, _ := x509.MarshalPKIXPublicKey(senderPub)
	receiverPub, _, _ := ed25519.GenerateKey(nil)
	receiverAddress, _ := x509.MarshalPKCS8PrivateKey(receiverPub)
	transaction := &blockchain.Transaction{
		Sender:   senderAddress,
		Receiver: receiverAddress,
		Amount:   10,
	}
	transaction.Sign(senderPriv)

	chain := blockchain.NewBlockchain()
	if err := chain.AddTransaction(transaction); err != nil {
		log.Fatalln("Invalid transaction")
	}
	chain.Mine()

}
