***Beth-go (Better Ethereum Library)***

Beth is a library that will provide functions that can interact with Ethereum consistently, given the issues that are observed in Ethereum. The major functions in this library (along with proposed code design) is discussed as follows:

**Transact**

This function will perform write operations (with retry capability) on the blockchain. `Transact` will also handle transaction details such as nonce and gas-price (based on the current fast gas-price estimate). However, the caller is expected to set gas-limit for each write transaction on the blockchain. This function will also wait for a user specified number of confirmation blocks to be seen on the blockchain post the executed transaction.


```
func Transact(ctx context.Context, preConditionCheck func() bool, f func() (tx, error), postConditionCheck func() bool, waitBlocks int) error {

    preconditionCheck(ctx)

    for {

        select {
        case <-ctx.Done():
            return ctx.Error()
        default:
        }

        func() {
            // Set gas-price in txOpts

            RetryNonceTx(innerCtx, f)

            if _, err := patchedWaitMined(tx); err == nil {
                return
            }

            // Retry tx if err is not nil

        }()

        if postCondition(ctx) {
            return nil
        }
    }
    txBlockNumber := TxBlockNumber(tx.Hash())

    currentBlockNumber := CurrentBlockNumber()

    while (currentBlockNumber - txBlockNumber) < waitBlocks {
        sleep()
        currentBlockNumber := CurrentBlockNumber()
    } 
}
```

**Transfer**

This function will transfer Eth to specified ethereum address using the previously defined `Transact` call.


**Get**

This function will execute read-only transactions on Ethereum until it does not return an error or until the given `context` is done. The values read from the operation are expected to be captured by the calling function.

**Balance**

This method will use the `Get` function to read the balances of the ethereum address.

```
func (client *Client) BalanceOf(ctx context.Context, addr common.Address, callOpts *bind.CallOpts) (val *big.Int, err error) {
    err = client.Get(
        ctx,
        callOpts,
        func(*bind.CallOpts) (err error) {
            val, err = client.ethClient.BalanceAt(ctx, addr, nil)
            return
        },
    )
    return
}
```
