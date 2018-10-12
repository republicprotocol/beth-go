package beth

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"log"
	"math"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
)

// ErrPreConditionCheckFailed indicates that the pre-condition for executing
// a transaction failed.
var ErrPreConditionCheckFailed = errors.New("pre-condition check failed")

// ErrPostConditionCheckFailed indicates that the post-condition for executing
// a transaction failed.
var ErrPostConditionCheckFailed = errors.New("post-condition check failed")

// ErrIncorrectNonce indicates that there exists another transaction with the
// same nonce and a higher or equal gas price as the present transaction.
var ErrIncorrectNonce = errors.New("there exists another transaction with the same nonce")

// Account is an Ethereum external account that can submit write transactions
// to the Ethereum blockchain. An Account is defined by its public address and
// respective private key.
type Account interface {
	EthClient() Client

	// Address returns the ethereum address of the account holder.
	Address() common.Address

	// Transfer sends the specified value of Eth to the given address.
	Transfer(ctx context.Context, to common.Address, value *big.Int, confirmBlocks int64) error

	// Transact performs a write operation on the Ethereum blockchain. It will
	// first conduct a preConditionCheck and if the check passes, it will
	// repeatedly execute the transaction followed by a postConditionCheck,
	// until the transaction passes and the postConditionCheck returns true.
	Transact(ctx context.Context, preConditionCheck func() bool, f func(bind.TransactOpts) (*types.Transaction, error), postConditionCheck func() bool, confirmBlocks int64) error
}

type account struct {
	mu     *sync.RWMutex
	client Client

	callOpts     bind.CallOpts
	transactOpts bind.TransactOpts
}

// NewAccount returns a user account for the provided private key which is
// connected to a ethereum client.
func NewAccount(url string, privateKey *ecdsa.PrivateKey) (Account, error) {
	client, err := Connect(url)
	if err != nil {
		return nil, err
	}

	transactOpts := *bind.NewKeyedTransactor(privateKey)

	// Retrieve nonce and update transactOpts.
	nonce, err := client.EthClient.PendingNonceAt(
		context.Background(),
		transactOpts.From)
	if err != nil {
		return nil, err
	}
	transactOpts.Nonce = big.NewInt(int64(nonce))

	ethAccount := &account{
		mu:     new(sync.RWMutex),
		client: client,

		callOpts:     bind.CallOpts{},
		transactOpts: transactOpts,
	}

	ethAccount.mu.Lock()
	defer ethAccount.mu.Unlock()

	// Retrieve and update transactOpts with the current 'fast' gas price
	ethAccount.updateGasPrice()

	return ethAccount, nil
}

// Address returns the ethereum address of the account.
func (account *account) Address() common.Address {
	return account.transactOpts.From
}

// EthClient returns the ethereum client that the account is connected to.
func (account *account) EthClient() Client {
	return account.client
}

// Transact attempts to execute a transaction on the Ethereum blockchain with
// the retry functionality.
func (account *account) Transact(ctx context.Context, preConditionCheck func() bool, f func(bind.TransactOpts) (*types.Transaction, error), postConditionCheck func() bool, waitForBlocks int64) error {

	// Do not proceed any further if the (not nil) pre-condition check fails
	if preConditionCheck != nil && !preConditionCheck() {
		return ErrPreConditionCheckFailed
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

			account.mu.Lock()
			defer account.mu.Unlock()

			account.updateGasPrice()

			// This will attempt to execute 'f' until no nonce error is
			// returned or if ctx times out
			tx, err := account.retryNonceTx(ctx, f)
			if err != nil {
				return err
			}

			receipt, err := account.client.WaitMined(ctx, tx)
			if err != nil {
				return err
			}

			// Transaction did not error, proceed to post-condition checks
			txHash = receipt.TxHash
			return nil
		}(); err != nil {
			// There is another transaction with the same nonce and a higher or
			// equal gas price as that of this transaction.
			if strings.Compare(err.Error(), ErrIncorrectNonce.Error()) == 0 {
				return err
			}
			log.Println(err)
			continue
		}

		// If post-condition check passes, proceed to wait for a specified
		// number of blocks to be confirmed after the transaction's block
		if postConditionCheck == nil || postConditionCheck() {
			break
		}

		// Wait for sometime before attempting to execute the transaction
		// again. If context is done, return error to indicate that
		// post-condition failed
		select {
		case <-ctx.Done():
			return ErrPostConditionCheckFailed
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
	blockNumber, err := account.client.TxBlockNumber(ctx, txHash.String())
	if err != nil {
		return err
	}

	// Attempt to get current block number. If context times out, an error will
	// be returned.
	currentBlockNumber, err := account.client.CurrentBlockNumber(ctx)
	if err != nil {
		return err
	}

	// Keep getting current block number until it is greater than
	// 'waitForBlocks' + transaction's block number. If context times out, the
	// error is returned.
	for big.NewInt(0).Sub(currentBlockNumber, blockNumber).Cmp(big.NewInt(waitForBlocks)) < 0 {
		currentBlockNumber, err = account.client.CurrentBlockNumber(ctx)
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

// Transfer transfers eth from the account to an ethereum address.
func (account *account) Transfer(ctx context.Context, to common.Address, value *big.Int, confirmBlocks int64) error {

	// Pre-condition check: Check if the account has enough balance
	preConditionCheck := func() bool {
		accountBalance, err := account.client.BalanceOf(ctx, account.Address(), &account.callOpts)
		return err == nil && accountBalance.Cmp(value) >= 0
	}

	// Transaction: Transfer eth to address
	f := func(transactOpts bind.TransactOpts) (*types.Transaction, error) {
		bound := bind.NewBoundContract(to, abi.ABI{}, nil, account.client.EthClient, nil)

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

// retryNonceTx retries transaction execution on the blockchain until nonce
// errors are not seen, or until the context times out.
func (account *account) retryNonceTx(ctx context.Context, f func(bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	tx, err := f(account.transactOpts)

	// On successful execution, increment nonce in transactOpts and return
	if err == nil {
		account.transactOpts.Nonce.Add(account.transactOpts.Nonce, big.NewInt(1))
		return tx, nil
	}

	// Process errors to check for nonce issues
	// If error indicates that nonce is too low, increment nonce and retry
	if err == core.ErrNonceTooLow || strings.Contains(err.Error(), "nonce is too low") {
		account.transactOpts.Nonce.Add(account.transactOpts.Nonce, big.NewInt(1))
		return account.retryNonceTx(ctx, f)
	}

	// If error indicates that nonce is too low, decrement nonce and retry
	if err == core.ErrNonceTooHigh {
		account.transactOpts.Nonce.Sub(account.transactOpts.Nonce, big.NewInt(1))
		return account.retryNonceTx(ctx, f)
	}

	// If any other type of nonce error occurs we will refresh the nonce and
	// try again for up to 1 minute
	var nonce uint64
	for try := 0; try < 60 && strings.Contains(err.Error(), "nonce"); try++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}

		// Get updated nonce and retry 'f'
		nonce, err = account.client.EthClient.PendingNonceAt(ctx, account.transactOpts.From)
		if err != nil {
			continue
		}
		account.transactOpts.Nonce = big.NewInt(int64(nonce))

		if tx, err = f(account.transactOpts); err == nil {
			account.transactOpts.Nonce.Add(account.transactOpts.Nonce, big.NewInt(1))
			return tx, nil
		}
	}

	if strings.Compare(err.Error(), core.ErrReplaceUnderpriced.Error()) == 0 {
		// Update nonce and return error
		account.transactOpts.Nonce.Add(account.transactOpts.Nonce, big.NewInt(1))
		return tx, ErrIncorrectNonce
	}

	return tx, err
}

// updateGasPrice will retrieve the current 'fast' gas price
// and update the account's transactOpts. This function expects the caller to
// handle potential data race conditions (i.e. Locking of mutex prior to
// calling this method)
func (account *account) updateGasPrice() {
	gasPrice := suggestedGasPrice()
	if gasPrice != nil {
		account.transactOpts.GasPrice = gasPrice
	}
}

// suggestedGasPrice returns the fast gas price that ethGasStation
// recommends for transactions to be mined on Ethereum blockchain.
func suggestedGasPrice() *big.Int {
	request, _ := http.NewRequest("GET", "https://ethgasstation.info/json/ethgasAPI.json", nil)

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("cannot connect to ethGasStationAPI: %v", err)
		return nil
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("received status code %v from ethGasStation", response.StatusCode)
		return nil
	}

	type resp struct {
		Fast float64 `json:"fast"`
	}

	data := new(resp)

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Printf("cannot decode json response from ethGasStation: %v", err)
		return nil
	}
	return big.NewInt(int64(data.Fast * math.Pow10(8)))
}
