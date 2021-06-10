# cryptocurrency

A simple cryptocurrency built out of curiosity. Not for real-world use.

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
- forks are not resolved (nodes just keep accepting valid blocks ad infinitum)
- peer-to-peer communication is unencrypted
- blockchain is not persisted to disk or compressed in any manner
- everything will probably implode if you actually run this in production

## Running a mining node

### Generating keys

Miners need to have an Ed25519 public key in order to receive coinbase transactions, i.e. the reward transaction for mining a block. To do that, you can generate a new random key pair first:

```shell
go run cmd/keygen/keygen.go
```

This will generate a key pair in your current directory, which will be used by default.

### Configuration

The node needs to be configured with a few environment variables if you want it to communicate with other nodes in the network:

```shell
export NODE_BIND_HOST=127.0.0.1:8000  # The address the HTTP server binds to
export NODE_ADVERTISED_HOST=some-public-ip:8000   # The address other nodes can reach your node with
```

Additionally you can define a seed host, which is used to synchronize the blockchain when the node initially starts:

```shell
export NODE_SEED_HOST=some-other-node:8080
```

### Compiling and running

Once you have your keys and you have configured the node, you can compile the app and start mining for blocks:

```shell
go build cmd/node/node.go
./node

...

2021/06/10 17:52:57 Adding genesis block: Block 0 000002be9afbfdaa977028a51d10bd590f9b56b03c3f570b8723e3809dc439ba transactions: 0
2021/06/10 17:52:57 Listening for API requests at localhost:8080
2021/06/10 17:53:13 🎉 Found valid block: Block 1 000002b6ba6b836c75c90bcbf938fafd75a0f5f933fd0b12fe18532477b2cc69 transactions: 1
```

### Peer-to-peer communication

If you want to run two nodes on the same host just for fun and to see how the peer-to-peer communication works, you can open two terminals (or use tmux or whatever) and run this in the first one:

```shell
export NODE_BIND_HOST=localhost:8000
./node
```

This will start the first node, which will immediately start mining for new blocks. Then you can start a second one using a different port. The second node will use the first one as a seed host, i.e. the second node will first synchronize any blocks the first node has mined before the two will start competing for the next block in the blockchain:

```shell
export NODE_BIND_HOST=localhost:8080
export NODE_SEED_HOST=localhost:8000
./node

...

2021/06/10 17:53:16 Adding genesis block: Block 0 000002be9afbfdaa977028a51d10bd590f9b56b03c3f570b8723e3809dc439ba transactions: 0
2021/06/10 17:53:16 Syncing blockchain via localhost:8080
2021/06/10 17:53:16 Listening for API requests at localhost:8000
2021/06/10 17:53:16 Remote node found valid block: Block 1 000002b6ba6b836c75c90bcbf938fafd75a0f5f933fd0b12fe18532477b2cc69 transactions: 1
2021/06/10 17:53:16 Node @ localhost:8080 sent greeting
2021/06/10 17:53:20 Remote node found valid block: Block 2 000003ef948b69e2167da552dade695ebda52cf0433964f96b10406b333ef333 transactions: 1
2021/06/10 17:53:37 🎉 Found valid block: Block 3 000003ad1c34b165cab51bf21d135bfc7783b51fca144aeb4b2204302c08292b transactions: 1
2021/06/10 17:53:42 Remote node found valid block: Block 4 0000021d3a4fff9de706de98a4dbc9233b9083fffbbf63db8f14adae7e15c20c transactions: 1
2021/06/10 17:53:57 Remote node found valid block: Block 5 0000077c4b9a8eea9d64a427175d7c26f2d115ee76358d7688e74e4edb4a49eb transactions: 1
2021/06/10 17:53:59 Remote node found valid block: Block 6 0000047082767c1f362c88bf0b98386c3a2d78f4c79bd66c20e511626c133141 transactions: 1
2021/06/10 17:54:06 🎉 Found valid block: Block 7 00000561aa27b01b5f709269cffb651a0fb7fe69aa31ee7b207e6067235b3b18 transactions: 1
2021/06/10 17:54:31 🎉 Found valid block: Block 8 000006dba152ad81e4f9ee69e05368ac63091ecd2173edb54b0154695abc859e transactions: 1
```

## Querying the blockchain

You can query a node for the state of the blockchain. For example, to get the latest block in the blockchain known by the node:

```shell
curl localhost:8080/api/v1/blockchain/ --silent | jq '.[-1]'
{
  "number": 16,
  "time": "2021-06-10T14:58:27.730607Z",
  "transactions": [
    {
      "sender": null,
      "receiver": "gBg426L2kNWAE1WFz+Jd+GmlcQ4XUabIqLvAAxz9OgI=",
      "amount": 10,
      "nonce": 0,
      "time": "2021-06-10T14:58:22.083545Z",
      "signature": null
    }
  ],
  "nonce": 640205,
  "previousHash": "AAAFkZjtFxv+jUoTTqrzcp8XC/9TmDPi2OM5VxpTgUA=",
  "hash": "AAADdOoIKvXfBAxx2RK9nIYIwOcgDfkjzEFXRycEHKc="
}
```

The transaction in the block above is a coinbase transaction, thus it has neither a valid signature nor a sender.
