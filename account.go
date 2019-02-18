package beth

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ErrPreConditionCheckFailed indicates that the pre-condition for executing
// a transaction failed.
var ErrPreConditionCheckFailed = errors.New("pre-condition check failed")

// ErrPostConditionCheckFailed indicates that the post-condition for executing
// a transaction failed.
var ErrPostConditionCheckFailed = errors.New("post-condition check failed")

// ErrNonceIsOutOfSync indicates that there exists another transaction with the
// same nonce and a higher or equal gas price as the present transaction.
var ErrNonceIsOutOfSync = errors.New("nonce is out of sync")

// ErrDuplicateAddress indicates that there exists another address that has
// been mapped to the same key.
var ErrDuplicateAddress = errors.New("the key has already been mapped to another address")

// ErrAddressNotFound indicates that the given key is not present in the
// address book.
var ErrAddressNotFound = errors.New("key does not have an entry in the address book")

// The TxExecutionSpeed indicates the tier of speed that the transaction falls
// under while writing to the blockchain.
type TxExecutionSpeed uint8

// TxExecutionSpeed values.
const (
	Nil = TxExecutionSpeed(iota)
	SafeLow
	Average
	Fast
	Fastest
)

// Account is an Ethereum external account that can submit write transactions
// to the Ethereum blockchain. An Account is defined by its public address and
// respective private key.
type Account interface {

	// Client that the account is connected to. Using the Client, read-only
	// operations can be executed on Ethereum.
	Client() Client

	// EthClient returns the actual Ethereum client.
	EthClient() *ethclient.Client

	// Address returns the ethereum address of the account holder.
	Address() common.Address

	// BalanceAt returns the wei balance of the account. The block number can be
	// nil, in which case the balance is taken from the latest known block.
	BalanceAt(ctx context.Context, blockNumber *big.Int) (*big.Int, error)

	// Store address in address book.
	WriteAddress(key string, address common.Address)

	// ReadAddress returns address mapped to the given key in the address book.
	ReadAddress(key string) (common.Address, error)

	// Transfer sends the specified value of Eth to the given address.
	Transfer(ctx context.Context, to common.Address, value, gasPrice *big.Int, confirmBlocks int64, sendAll bool) (*types.Transaction, error)

	// Transact performs a write operation on the Ethereum blockchain. It will
	// first conduct a preConditionCheck and if the check passes, it will
	// repeatedly execute the transaction followed by a postConditionCheck,
	// until the transaction passes and the postConditionCheck returns true.
	// Transact will immediately stop retrying if an ErrReplaceUnderpriced is
	// returned from ethereum.
	Transact(ctx context.Context, preConditionCheck func() bool, f func(*bind.TransactOpts) (*types.Transaction, error), postConditionCheck func() bool, confirmBlocks int64) (*types.Transaction, error)

	// Sign the given message with the account's private key.
	Sign(msgHash []byte) ([]byte, error)

	// SetGasPrice allows the account holder to set the gasPrice to a specific
	// value.
	SetGasPrice(gasPrice float64)

	// ResetToPendingNonce will wait for a 'coolDown' time (in milliseconds)
	// before updating transaction nonce to current pending nonce.
	ResetToPendingNonce(ctx context.Context, coolDown time.Duration) error

	// FormatTransactionView returns the formatted string with the URL at which
	// the transaction can be viewed.
	FormatTransactionView(msg, txHash string) (string, error)

	NewERC20(addressOrAlias string) (ERC20, error)
}

type account struct {
	mu     *sync.RWMutex
	client Client

	callOpts     *bind.CallOpts
	transactOpts *bind.TransactOpts

	privateKey *ecdsa.PrivateKey

	addressBook AddressBook
}

// NewAccount returns a user account for the provided private key which is
// connected to an Ethereum client.
func NewAccount(url string, privateKey *ecdsa.PrivateKey) (Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Connect to client
	client, err := Connect(url)
	if err != nil {
		return nil, err
	}

	// Setup transact opts
	transactOpts := bind.NewKeyedTransactor(privateKey)
	nonce, err := client.EthClient().PendingNonceAt(ctx, transactOpts.From)
	if err != nil {
		return nil, err
	}
	transactOpts.Nonce = big.NewInt(0).SetUint64(nonce)

	netID, err := client.ethClient.NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	// Create account
	account := &account{
		mu:     new(sync.RWMutex),
		client: client,

		callOpts:     new(bind.CallOpts),
		transactOpts: transactOpts,

		privateKey: privateKey,

		addressBook: DefaultAddressBook(netID.Int64()),
	}

	return account, nil
}

// Address returns the ethereum address of the account.
func (account *account) Address() common.Address {
	return account.transactOpts.From
}

// BalanceAt returns the wei balance of the account. The block number can be nil,
// in which case the balance is taken from the latest known block.
func (account *account) BalanceAt(ctx context.Context, blockNumber *big.Int) (*big.Int, error) {
	return account.client.ethClient.BalanceAt(ctx, account.Address(), nil)
}

// WriteAddress to the address book, overwrite if already exists
func (account *account) WriteAddress(key string, address common.Address) {
	account.addressBook[key] = address
}

// ReadAddress from the address book, return an error if the address does not
// exist
func (account *account) ReadAddress(key string) (common.Address, error) {
	if address, ok := account.addressBook[key]; ok {
		return address, nil
	}
	return common.Address{}, ErrAddressNotFound
}

func (account *account) Client() Client {
	return account.client
}

func (account *account) EthClient() *ethclient.Client {
	return account.client.EthClient()
}

// Transact attempts to execute a transaction on the Ethereum blockchain with
// the retry functionality. It stops retrying if tx is completed without any
// error, or if given context times-out, or if ErrReplaceUnderpriced is
// returned from Ethereum.
func (account *account) Transact(ctx context.Context, preConditionCheck func() bool, f func(*bind.TransactOpts) (*types.Transaction, error), postConditionCheck func() bool, waitForBlocks int64) (*types.Transaction, error) {

	// Do not proceed any further if the (not nil) pre-condition check fails
	if preConditionCheck != nil && !preConditionCheck() {
		return nil, ErrPreConditionCheckFailed
	}

	sleepDurationMs := time.Duration(1000)
	var txHash common.Hash
	var transction *types.Transaction

	// Keep retrying 'f' until the post-condition check passes or the context
	// times out.
	var postConPassed = false
	for !postConPassed {
		// If context is done, return error
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if err := func() error {
			account.mu.Lock()
			defer account.mu.Unlock()

			account.updateGasPrice(Fast)
			// This will attempt to execute 'f' until no nonce error is
			// returned or if ctx times out
			innerCtx, innerCancel := context.WithTimeout(ctx, 10 * time.Minute)
			defer innerCancel()

			tx, err := account.retryNonceTx(innerCtx, f)
			if err != nil {
				return err
			}

			if _, err := account.client.WaitMined(innerCtx, tx); err != nil {
				return err
			}
			txHash = tx.Hash()
			transction = tx

			// Transaction did not error, proceed to post-condition checks
			return nil
		}(); err != nil {
			// There is another transaction with the same nonce and a higher or
			// equal gas price as that of this transaction.
			if strings.Compare(err.Error(), core.ErrReplaceUnderpriced.Error()) == 0 {
				return nil, ErrNonceIsOutOfSync
			}
			fmt.Println(err)
		}

		for i := 0; i < 180; i++ {
			select {
			case <-ctx.Done():
				return nil, ErrPostConditionCheckFailed
			default:
				if postConditionCheck == nil || postConditionCheck() {
					postConPassed = true
				}
			}
			if postConPassed {
				break
			}
			time.Sleep(time.Second)
		}

		// If post-condition check passes, proceed to wait for a specified
		// number of blocks to be confirmed after the transaction's block

		// Wait for sometime before attempting to execute the transaction
		// again. If context is done, return error to indicate that
		// post-condition failed
		select {
		case <-ctx.Done():
			return nil, ErrPostConditionCheckFailed
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
		return nil, err
	}

	// Attempt to get current block number. If context times out, an error will
	// be returned.
	currentBlockNumber, err := account.client.CurrentBlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	// Keep getting current block number until it is greater than
	// 'waitForBlocks' + transaction's block number. If context times out, the
	// error is returned.
	for big.NewInt(0).Sub(currentBlockNumber, blockNumber).Cmp(big.NewInt(waitForBlocks)) < 0 {
		currentBlockNumber, err = account.client.CurrentBlockNumber(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Second):
			}
			continue
		}
	}
	return transction, nil
}

// Transfer transfers eth from the account to an ethereum address. If the value
// is nil then it transfers all the balance to the `to` address.
func (account *account) Transfer(ctx context.Context, to common.Address, value, gasPrice *big.Int, confirmBlocks int64, sendAll bool) (*types.Transaction, error) {
	if value == nil {
		return nil, fmt.Errorf("value cannot be nil")
	}

	// Pre-condition check: Check if the account has enough balance
	preConditionCheck := func() bool {
		if value == nil {
			return true
		}
		accountBalance, err := account.client.BalanceOf(ctx, account.Address())
		return err == nil && accountBalance.Cmp(value) >= 0
	}

	// Transaction: Transfer eth to address
	f := func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		bound := bind.NewBoundContract(to, abi.ABI{}, nil, account.client.EthClient(), nil)
		if gasPrice == nil {
			gasPrice = transactOpts.GasPrice
		}

		if sendAll {
			balance, err := account.BalanceAt(ctx, nil)
			if err != nil {
				return nil, err
			}
			value = new(big.Int).Sub(balance, new(big.Int).Mul(big.NewInt(21000), gasPrice))
		}

		transactor := &bind.TransactOpts{
			From:     transactOpts.From,
			Signer:   transactOpts.Signer,
			Value:    value,
			GasLimit: 21000,
			GasPrice: gasPrice,
			Context:  ctx,
		}
		if transactOpts.Nonce != nil {
			transactor.Nonce = big.NewInt(0).Set(transactOpts.Nonce)
		}
		if transactor.GasPrice == nil {
			transactor.GasPrice = big.NewInt(0).Set(transactOpts.GasPrice)
		}

		tx, err := bound.Transfer(transactor)
		if err != nil {
			return tx, err
		}
		return tx, nil
	}
	return account.Transact(ctx, preConditionCheck, f, nil, confirmBlocks)
}

// Sign the given message with the account's private key.
func (account *account) Sign(msgHash []byte) ([]byte, error) {
	return crypto.Sign(msgHash, account.privateKey)
}

// SetGasPrice will allow the caller to set gas price of transactOpts.
func (account *account) SetGasPrice(gasPrice float64) {
	account.mu.Lock()
	defer account.mu.Unlock()

	account.transactOpts.GasPrice = big.NewInt(int64(gasPrice * math.Pow10(9)))
}

// ResetToPendingNonce will allow the caller to reset nonce to pending nonce.
// This function will wait for a 'coolDown' time (in milliseconds) before
// updating the nonce in transactOpts.
func (account *account) ResetToPendingNonce(ctx context.Context, coolDown time.Duration) error {
	account.mu.Lock()
	defer account.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(coolDown * time.Millisecond):
	}

	// Get pending nonce
	nonce, err := account.client.EthClient().PendingNonceAt(ctx, account.transactOpts.From)
	if err != nil {
		return err
	}
	account.transactOpts.Nonce = big.NewInt(int64(nonce))
	return nil
}

// FormatTransactionView returns the formatted string with the URL at which the
// transaction can be viewed.
func (account *account) FormatTransactionView(msg, txHash string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	netID, err := account.client.ethClient.NetworkID(ctx)
	if err != nil {
		return "", err
	}
	switch netID.Int64() {
	case 1:
		return fmt.Sprintf("%s, the transaction can be viewed at https://etherscan.io/tx/%s", msg, txHash), nil
	case 3:
		return fmt.Sprintf("%s, the transaction can be viewed at https://ropsten.etherscan.io/tx/%s", msg, txHash), nil
	case 42:
		return fmt.Sprintf("%s, the transaction can be viewed at https://kovan.etherscan.io/tx/%s", msg, txHash), nil
	default:
		return "", fmt.Errorf("unknown network id : %d", netID.Int64())
	}
}

// SuggestedGasPrice returns the gas price that ethGasStation recommends for
// transactions to be mined on Ethereum blockchain based on the speed provided.
func SuggestedGasPrice(txSpeed TxExecutionSpeed) (*big.Int, error) {
	request, err := http.NewRequest("GET", "https://ethgasstation.info/json/ethgasAPI.json", nil)
	if err != nil {
		return nil, fmt.Errorf("cannot build request to ethGasStation = %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	res, err := (&http.Client{}).Do(request)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to ethGasStationAPI = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %v from ethGasStation", res.StatusCode)
	}

	data := struct {
		SafeLow float64 `json:"safeLow"`
		Average float64 `json:"average"`
		Fast    float64 `json:"fast"`
		Fastest float64 `json:"fastest"`
	}{}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("cannot decode response body from ethGasStation = %v", err)
	}

	switch txSpeed {
	case SafeLow:
		return big.NewInt(int64(data.SafeLow * math.Pow10(8))), nil
	case Average:
		return big.NewInt(int64(data.Average * math.Pow10(8))), nil
	case Fast:
		return big.NewInt(int64(data.Fast * math.Pow10(8))), nil
	case Fastest:
		return big.NewInt(int64(data.Fastest * math.Pow10(8))), nil
	default:
		return nil, fmt.Errorf("invalid speed tier: %v", txSpeed)
	}
}

// retryNonceTx retries transaction execution on the blockchain until nonce
// errors are not seen, or until the context times out.
func (account *account) retryNonceTx(ctx context.Context, f func(*bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	transactor := &bind.TransactOpts{
		From:     account.transactOpts.From,
		Signer:   account.transactOpts.Signer,
		Value:    big.NewInt(0),
		GasLimit: account.transactOpts.GasLimit,
		Context:  ctx,
	}
	if account.transactOpts.Nonce != nil {
		transactor.Nonce = big.NewInt(0).Set(account.transactOpts.Nonce)
	}
	if account.transactOpts.GasPrice != nil {
		transactor.GasPrice = big.NewInt(0).Set(account.transactOpts.GasPrice)
	}

	tx, err := f(transactor)

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
	if err == core.ErrNonceTooHigh || strings.Contains(err.Error(), "nonce is too high") {
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
		nonce, err = account.client.EthClient().PendingNonceAt(ctx, account.transactOpts.From)
		if err != nil {
			continue
		}
		account.transactOpts.Nonce = big.NewInt(int64(nonce))

		if tx, err = f(account.transactOpts); err == nil {
			account.transactOpts.Nonce.Add(account.transactOpts.Nonce, big.NewInt(1))
			return tx, nil
		}
	}

	return tx, err
}

// updateGasPrice will retrieve the current 'fast' gas price
// and update the account's transactOpts. This function expects the caller to
// handle potential data race conditions (i.e. Locking of mutex prior to
// calling this method)
func (account *account) updateGasPrice(txSpeed TxExecutionSpeed) error {
	gasPrice, err := SuggestedGasPrice(txSpeed)
	if err != nil {
		return err
	}
	if gasPrice != nil {
		account.transactOpts.GasPrice = gasPrice
	}
	return nil
}
