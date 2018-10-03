***Design Document***


**`SendTx(f func() (tx,error)) (tx, error)`**

This method will execute the function `f` until no nonce related errors are observed. It will finally return a transaction and error. _(Same as the implementation in `republic-go`)_. 


**`SuggestGasPrice(ctx) (uint64, error)`**

- go-ethereum : `ethClient.SuggestGasPrice(ctx)` retrieves the currently suggested gas price to allow a timely execution of a transaction.
- web3 : `web3.eth.gasPrice` returns the current gas price, which is determined by the x latest blocks median gas price. 

_Note: Not sure if these methods are completely reliable since `web3` gets its value from the nodes. See [issue](https://github.com/ethereum/web3.js/issues/1230)._

_Update: The gasPrice method will be invoked until a 2/3rds majority answer is received. If a majority is not received even after a specified time, the average of all the values received uptil that point is returned._ 


**`Write(context, preCondition func(), f func(), postCondition func()) (error)`**

This method will execute the following functions (in the order it appears).

- `preCondition func() bool` will execute the necessary checks before executing the transaction on the blockchain. It will return true, if the checks pass. Otherwise, it will retry (with exponential delay on each attempt) until the checks pass or the context is done.

- `SuggestGasPrice` to set gas price for the transaction to be executed.

- `f func() (tx, error)` is executed once the `preCondition` returns true. This will attempt to write the transaction on the blockchain.

- `postCondition func() bool` : Rather than waiting for `patchedWaitMined`, a `postCondition` function which will check the state of the write operation will be executed (i.e, For a `settle` operation, the status of the order is checked to see if it is `Settled` on the blockchain). This function will wait for an approximate number of blocks before checking the status and `postCondition` will be retried once if the check does not pass. 

`Write` will retry from `SuggestGasPrice` (with exponential time delay, starting with 1 second) until `ctx` is done or `postCondition` returns true.

```
func Write(ctx context.Context, preConditionCheck func() bool, f func() (tx, error), postConditionCheck func() bool) error {

    preconditionCheck(ctx)

    for {

        select {
        case <-ctx.Done():
            return ctx.Error()
        default:
        }

        func() {
            innerCtx, innerCancel := context.WithTimeout(ctx, time.Minute)
            defer innerCancel()

            // TODO: Get suggestion gas price.

            gasLimitEstimates := [3]int{}
            co.ParForAll(3, func(i int) {
                // TODO: Estimate gas limit.
                gasLimitEstimates[i] = // ...
            })

            RetryNonceTx(innerCtx, f)

            if _, err := patchedWaitMined(tx); err == nil {
                return
            }

            // Transaction failed, sleep for 6 blocks before checking 
            // postCondition and retrying.
            sleep(for approx 6 blocks)

        }()

        if postCondition(ctx) {
            return nil
        }
    }
}
```


**`Read(ctx, f func()) (interface{}, error)`**

Read will execute the function `f` until it does not return an error or until the `ctx` is done.


**`TransferEth(transactOpts, value) (error)`**

This function will be used to transfer eth from one ethereum address to another.

```
if balance (from) ≥ value + tx_fees {

    gasPrice = SuggestGasPrice(ctx)
    Set gasPrice in transactOpts

    tx, err = SendTx( func() (tx, error) {
        return TransferEth(transactOpts)
    })

    if err != nil {
        return err
    }

    Call patchedWaitMined (tx) and return err
}

return Error("Insufficient funds")
```


**`BalanceOf(token) (value, error)`**

This method will use the earlier defined `Read` function and return the result of `Read`.

```
Read(ctx, func() (value, error) {
    return ethereum.BalanceOf(token)
})
```
