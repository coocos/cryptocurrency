# cryptocurrency

A simple cryptocurrency built out of curiosity.

## Features

- classic blockchain structure
- parallel SHA256-based proof-of-work computation
- signed transactions using Ed25519
- miners are rewarded with a coinbase transaction per block
- balance based account model
- peer-to-peer networking on top of HTTP

## Limitations

- no transaction rewards (miners have no incentive to include transactions in blocks)
- no difficulty scaling (increasing amount of mining nodes will lead to rapid inflation)
- transactions vulnerable to replay attacks
- peer-to-peer communication is unencrypted

## Running

Miners need to have an Ed25519 public key in order to receive coinbase transactions, i.e. the reward transaction for mining a block. To do that, you can generate a new random key pair first:

```shell
go run cmd/keygen/keygen.go
```

This will generate a key pair in your current directory, which will be used by default. After that, you can start mining:

```shell
go run cmd/node/node.go
```
