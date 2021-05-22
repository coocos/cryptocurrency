package blockchain

// EventEmitter is an interface for broadcasting blockhain change events
type EventEmitter interface {
	EmitBlock(Block)
	EmitTransaction(Transaction)
}

// DummyEventEmitter swallows received blocks and transactions
type DummyEventEmitter struct {
}

func (b *DummyEventEmitter) EmitBlock(block Block)                   {}
func (b *DummyEventEmitter) EmitTransaction(transaction Transaction) {}
