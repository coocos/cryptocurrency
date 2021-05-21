# cryptocurrency

A simple cryptocurrency built out of curiosity.

## Features

- classic blockchain structure
- parallel SHA256-based proof-of-work computation
- signed transactions using [Ed25519](https://en.wikipedia.org/wiki/EdDSA#Ed25519)
- miners are rewarded with a coinbase transaction per block
- balance based account model

## Not implemented (yet)

- peer-to-peer networking & consensus (i.e. the blockchain is currenly mined by a single node)
- transaction rewards
- difficulty scaling

## Running

Miners need to have a valid [Ed25519](https://en.wikipedia.org/wiki/EdDSA#Ed25519) key pair in order to receive coinbase transactions, i.e. the reward transaction for mining a block. To do that, you can generate a new random key pair first:

```shell
go run cmd/keygen.go
```

This will generate a key pair in your current directory, which will be used by default. Once you have your key pair generated, you can start mining:

```shell
go run cmd/node/main.go
```
