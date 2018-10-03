package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/eth-go/utils"
)

// ErrorPreConditionCheckFailed indicates that the pre-condition for executing
// a transaction failed.
var ErrorPreConditionCheckFailed = errors.New("pre-condition check failed")

// ErrorPostConditionCheckFailed indicates that the post-condition for executing
// a transaction failed.
var ErrorPostConditionCheckFailed = errors.New("post-condition check failed")

// Account is the Ethereum account that can perform read and write transactions
// on the Ethereum blockchain.
type Account struct {
	mu     *sync.RWMutex
	client Client

	callOpts     bind.CallOpts
	transactOpts bind.TransactOpts
}

// NewAccount returns a user account for the provided private key which is
// connected to a ethereum client.
func NewAccount(url string, privateKey *ecdsa.PrivateKey) (*Account, error) {
	client, err := Connect(url)
	if err != nil {
		return nil, err
	}

	transactOpts := *bind.NewKeyedTransactor(privateKey)

	// Retrieve nonce and update transactOpts.
	nonce, err := client.ethClient.PendingNonceAt(
		context.Background(),
		transactOpts.From)
	if err != nil {
		return nil, err
	}
	transactOpts.Nonce = big.NewInt(int64(nonce))

	account := &Account{
		mu:     new(sync.RWMutex),
		client: client,

		callOpts:     bind.CallOpts{},
		transactOpts: transactOpts,
	}

	account.mu.Lock()
	defer account.mu.Unlock()

	// Retrieve and update transactOpts with the current 'fast' gas price
	account.updateGasPrice()

	return account, nil
}

// Address returns the ethereum address of the account.
func (account *Account) Address() common.Address {
	return account.transactOpts.From
}

// Transfer transfers eth from the account to an ethereum address.
func (account *Account) Transfer(
	ctx context.Context,
	to common.Address,
	value *big.Int,
	confirmBlocks int64,
) error {

	// Pre-condition check: Check if the account has enough balance
	preConditionCheck := func() bool {
		accountBalance, err := account.client.BalanceOf(
			ctx,
			account.Address(),
			&account.callOpts,
		)
		if err != nil || accountBalance.Cmp(value) <= 0 {
			return false
		}

		return true
	}

	// Transaction: Transfer eth to address
	f := func(transactOpts bind.TransactOpts) (*types.Transaction, error) {
		bound := bind.NewBoundContract(
			to,
			abi.ABI{},
			nil,
			account.client.ethClient,
			nil,
		)

		transactor := &bind.TransactOpts{
			From:     transactOpts.From,
			Nonce:    transactOpts.Nonce,
			Signer:   transactOpts.Signer,
			Value:    value,
			GasPrice: transactOpts.GasPrice,
			GasLimit: 21000,
			Context:  transactOpts.Context,
		}
		return bound.Transfer(transactor)
	}

	return account.Transact(ctx, preConditionCheck, f, nil, confirmBlocks)
}

// Transact attempts to execute a transaction on the Ethereum blockchain with
// the retry functionality.
func (account *Account) Transact(
	ctx context.Context,
	preConditionCheck func() bool,
	f func(bind.TransactOpts) (*types.Transaction, error),
	postConditionCheck func(ctx context.Context) bool,
	waitForBlocks int64,
) error {

	// Do not proceed any further if the (not nil) pre-condition check fails
	if preConditionCheck != nil && !preConditionCheck() {
		return ErrorPreConditionCheckFailed
	}

	sleepDurationMs := time.Duration(1000)
	var txHash common.Hash

	// Keep retrying 'f' until the post-condition check passes or the context
	// times out.
	for {
		// If context is done, return error
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := func() error {
			var err error

			innerCtx, innerCancel := context.WithTimeout(ctx, time.Minute)
			defer innerCancel()

			account.mu.Lock()
			defer account.mu.Unlock()

			account.updateGasPrice()

			// This will attempt to execute 'f' until no nonce error is
			// returned or if innerCtx times out
			tx, err := account.retryNonceTx(innerCtx, f)
			if err != nil {
				return err
			}

			receipt, err := account.client.WaitMined(innerCtx, tx)
			if err != nil {
				return err
			}

			// Transaction did not error, proceed to post-condition checks
			txHash = receipt.TxHash
			return nil
		}(); err != nil {
			log.Println(err)
			continue
		}

		// If post-condition check passes, proceed to wait for a specified
		// number of blocks to be confirmed after the transaction's block
		if postConditionCheck == nil || postConditionCheck(ctx) {
			break
		}

		// Wait for sometime before attempting to execute the transaction
		// again. If context is done, return error to indicate that
		// post-condition failed
		select {
		case <-ctx.Done():
			return ErrorPostConditionCheckFailed
		case <-time.After(sleepDurationMs * time.Millisecond):

		}

		// Increase delay for next round but saturate at 30s
		sleepDurationMs = time.Duration(float64(sleepDurationMs) * 1.6)
		if sleepDurationMs > 30000 {
			sleepDurationMs = 30000
		}
	}

	// At this point the transaction executed successfully and the
	// post-condition check has passed. This means that we can now proceed to
	// wait for a pre-defined number of blocks to be confirmed on the
	// blockchain after the transaction's block is confirmed

	// Attempt to get block number of transaction. If context times out, an
	// error will be returned.
	blockNumber, err := account.client.GetBlockNumberByTxHash(
		ctx,
		txHash.String(),
	)
	if err != nil {
		return err
	}

	// Attempt to get current block number. If context times out, an error will
	// be returned.
	currentBlockNumber, err := account.client.GetCurrentBlockNumber(ctx)
	if err != nil {
		return err
	}

	// Keep getting current block number until it is greater than
	// 'waitForBlocks' + transaction's block number. If context times out, the
	// error is returned.
	for big.NewInt(0).Sub(currentBlockNumber, blockNumber).Cmp(big.NewInt(waitForBlocks)) < 0 {
		currentBlockNumber, err = account.client.GetCurrentBlockNumber(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Millisecond):
			}
			continue
		}
	}
	return nil
}

// retryNonceTx retries transaction execution on the blockchain until nonce
// errors are not seen, or until the context times out.
func (account *Account) retryNonceTx(
	ctx context.Context,
	f func(bind.TransactOpts) (*types.Transaction, error),
) (*types.Transaction, error) {

	tx, err := f(account.transactOpts)

	// On successful execution, increment nonce in transactOpts and return
	if err == nil {
		account.transactOpts.Nonce.Add(account.transactOpts.Nonce,
			big.NewInt(1))
		return tx, nil
	}

	// Process errors to check for nonce issues

	// If error indicates that nonce is too low, increment nonce and retry
	if err == core.ErrNonceTooLow ||
		err == core.ErrReplaceUnderpriced ||
		strings.Contains(err.Error(), "nonce is too low") {

		account.transactOpts.Nonce.Add(account.transactOpts.Nonce,
			big.NewInt(1))
		return account.retryNonceTx(ctx, f)
	}

	// If error indicates that nonce is too low, decrement nonce and retry
	if err == core.ErrNonceTooHigh {
		account.transactOpts.Nonce.Sub(account.transactOpts.Nonce,
			big.NewInt(1))
		return account.retryNonceTx(ctx, f)
	}

	// If any other type of nonce error occurs we will refresh the nonce and
	// try again for up to 1 minute
	var nonce uint64
	for try := 0; try < 60 && strings.Contains(err.Error(), "nonce"); try++ {

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1 * time.Second):
		}

		// Get updated nonce and retry 'f'
		nonce, err = account.client.ethClient.PendingNonceAt(ctx,
			account.transactOpts.From)
		if err != nil {
			continue
		}
		account.transactOpts.Nonce = big.NewInt(int64(nonce))

		if tx, err = f(account.transactOpts); err == nil {
			account.transactOpts.Nonce.Add(account.transactOpts.Nonce,
				big.NewInt(1))
			return tx, nil
		}
	}
	return tx, err
}

// updateGasPrice will retrieve the current 'fast' gas price
// and update the account's transactOpts. This function expects the caller to
// handle potential data race conditions (i.e. Locking of mutex prior to
// calling this method)
func (account *Account) updateGasPrice() {
	gasPrice := utils.SuggestedGasPrice()
	if gasPrice != nil {
		account.transactOpts.GasPrice = gasPrice
	}
}
