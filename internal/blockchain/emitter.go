package blockchain

// EventEmitter is an interface for broadcasting blockhain change events
type EventEmitter interface {
	EmitBlock(Block)
	EmitTransaction(Transaction)
}

// DummyEventEmitter swallows received blocks and transactions
type DummyEventEmitter struct {
}

// EmitBlock simply swallows the block and does nothing
func (b *DummyEventEmitter) EmitBlock(block Block) {}

// EmitBlock simply swallows the block and does nothing
func (b *DummyEventEmitter) EmitTransaction(transaction Transaction) {}
