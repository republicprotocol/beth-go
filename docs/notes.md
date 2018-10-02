***Notes***

1. Currently, Parity nodes return the following error when a call to get current block number is issued:

`missing required field 'mixHash' for Header`

This affects the feature of waiting for a fixed number of blocks on the blockchain to be confirmed while executing a `Transact` call.

[go-ethereum issue on GitHub](https://github.com/ethereum/go-ethereum/issues/3230)

2. `TransactionByHash` does not return the block information of the transaction (although this is available). 

[go-ethereum issue on GitHub](https://github.com/ethereum/go-ethereum/issues/15210)

[parity-ethereum issue on GitHub](https://github.com/paritytech/parity-ethereum/issues/8841)

There is a [PR](https://github.com/ethereum/go-ethereum/pull/17662) that is a potential fix for this issue. 

**Current implementation**

Interact with infura's api to get the current block number and block details of any given transaction.
