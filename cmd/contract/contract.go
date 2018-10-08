package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	eth "github.com/republicprotocol/eth-go"
	"github.com/republicprotocol/eth-go/cmd/contract/bindings"
	"github.com/republicprotocol/republic-go/crypto"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Parse CLI args and get new Account and contract
	account, contract, err := newAccount()
	if err != nil {
		fmt.Printf("\n[error] %v\n", err)
		return
	}

	// Perform list and map operations in the contract
	listOperations(account, contract, 2)

	// Perform integer operations: set, read, increment
	integerOperations(account, contract)
}

func newAccount() (eth.UserAccount, *bindings.Bethtest, error) {

	// Check if all expected values were provided
	if len(os.Args) != 4 {
		return nil, nil, errors.New("\nInvalid number of arguments!\x1b[0m \n\nPlease enter a \x1b[37;1mkeystore path\x1b[0m, a \x1b[37;1mpassphrase\x1b[0m to unlock the keystore and a \x1b[37;1minfura network\x1b[0m\n\n\x1b[31;1m[Usage] go run contract_test.go <path/to/keystore/file> <passphrase> <kovan/ropsten>\x1b[0m")
	}

	// Parse command-line arguments
	keystorePath := os.Args[1]
	passphrase := os.Args[2]
	network := os.Args[3]

	// Ropsten : 0x46bcff69b2d5a677c40c05c1f034ef7bf0ee4742
	// Kovan : 0x055d30956deea82bfe0f99a2771dcc36a18dc9bb
	contractAddr := common.Address{}
	switch network {
	case "kovan":
		contractAddr = common.HexToAddress("0x055d30956deea82bfe0f99a2771dcc36a18dc9bb")
	case "ropsten":
		contractAddr = common.HexToAddress("0x46bcff69b2d5a677c40c05c1f034ef7bf0ee4742")
	default:
		return nil, nil, errors.New("invalid infura network")
	}

	// Decrypt keystore
	ks := crypto.Keystore{}

	// Open keystore file
	keyin, err := os.Open(keystorePath)
	if err != nil {
		return nil, nil, err
	}
	json, err := ioutil.ReadAll(keyin)
	if err != nil {
		return nil, nil, err
	}

	var privKey *ecdsa.PrivateKey

	// Decrypt private key using keystore and passphrase
	if err := ks.DecryptFromJSON(json, passphrase); err != nil {
		key, err := keystore.DecryptKey(json, passphrase)
		if err != nil {
			return nil, nil, err
		}
		privKey = key.PrivateKey
	} else {
		privKey = ks.EcdsaKey.PrivateKey
	}

	// Return a user account to perform transactions
	account, err := eth.NewUserAccount(fmt.Sprintf("https://%s.infura.io", network), privKey)
	if err != nil {
		return nil, nil, err
	}

	// Get contract
	contract, err := bindings.NewBethtest(contractAddr, bind.ContractBackend(account.EthClient().EthClient))
	if err != nil {
		return nil, nil, err
	}

	return account, contract, nil
}

func listOperations(account eth.UserAccount, contract *bindings.Bethtest, n int) {
	// Append random 'n' values to the list
	values := appendToList(n, contract, account)

	// Delete these values from the list
	deleteFromList(values, n, contract, account)
}

func integerOperations(account eth.UserAccount, contract *bindings.Bethtest) {
	// Context for 2 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Generate random value
	val := big.NewInt(int64(rand.Intn(100)))

	// Set value in the contract
	setInt(ctx, account, contract, val)

	// Increment the value in the contract
	increment(ctx, account, contract, val)
}

func appendToList(n int, contract *bindings.Bethtest, account eth.UserAccount) []*big.Int {
	// Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(n+2)*time.Minute)
	defer cancel()

	var val *big.Int
	var err error

	values := []*big.Int{}

	for i := 0; i < n; i++ {

		// Randomly create a value and append it to list
		val = big.NewInt(int64(rand.Intn(100)))
		values = append(values, val)

		fmt.Printf("\n\x1b[37;1mAppending %v to list\x1b[0m", val.String())

		// Pre-condition: Does element already exist in the list?
		preCondition := func(ctx context.Context) bool {
			exists := elementExists(ctx, account.EthClient(), contract, val)
			if exists {
				fmt.Println("\n[warning] Element exists in list!")
			}

			return !exists
		}

		// Append to list
		f := func(txOpts bind.TransactOpts) (*types.Transaction, error) {
			return contract.Append(&txOpts, val)
		}

		// Post-condition: Has element been added to the list?
		postCondition := func(ctx context.Context) bool {
			return elementExists(ctx, account.EthClient(), contract, val)
		}

		// Execute transaction
		account.Transact(ctx, preCondition, f, postCondition, 1)
	}

	// Print the size of the new list
	if val, err = size(ctx, account.EthClient(), contract); err == nil {
		fmt.Printf("\n\x1b[37;1mThe new size of array is %v\n\x1b[0m", val.String())
	}

	return values
}

func deleteFromList(values []*big.Int, n int, contract *bindings.Bethtest, account eth.UserAccount) {
	// Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(n+2)*time.Minute)
	defer cancel()

	var err error

	for _, val := range values {
		fmt.Printf("\n\x1b[37;1mDeleting %v from list\x1b[0m", val.String())

		// Pre-condition: is list is empty or is element absent in list?
		preCondition := func(ctx context.Context) bool {

			// Read size of list
			size, err := size(ctx, account.EthClient(), contract)
			if err != nil || size.Cmp(big.NewInt(0)) <= 0 {
				fmt.Println("\n[warning] list is empty!")
				return false
			}

			// Check if element is present
			exists := elementExists(ctx, account.EthClient(), contract, val)
			if !exists {
				fmt.Printf("\n[warning] %v is not present in list!\n", val.String())
			}

			return exists
		}

		// Remove element
		f := func(txOpts bind.TransactOpts) (*types.Transaction, error) {
			return contract.Remove(&txOpts, val)
		}

		// Post-condition: Has element been deleted successfully?
		postCondition := func(ctx context.Context) bool {
			return !elementExists(ctx, account.EthClient(), contract, val)
		}

		// Execute delete tx
		account.Transact(ctx, preCondition, f, postCondition, 1)
	}

	// Print size of new list
	val, err := size(ctx, account.EthClient(), contract)
	if err == nil {
		fmt.Printf("\n\x1b[37;1mThe size of array after deleting elements is %v\n\x1b[0m", val.String())
	}
}

func setInt(ctx context.Context, account eth.UserAccount, contract *bindings.Bethtest, val *big.Int) error {

	fmt.Printf("\n\x1b[37;1mSetting integer %v in the contract\n\x1b[0m", val.String())

	// Set integer in contract
	f := func(txOpts bind.TransactOpts) (*types.Transaction, error) {
		return contract.Set(&txOpts, val)
	}

	// Post-condition: Confirm that the integer has the new value
	postCondition := func(ctx context.Context) bool {
		newVal, err := read(ctx, account.EthClient(), contract)
		if err != nil {
			return false
		}
		return newVal.Cmp(val) == 0
	}

	return account.Transact(ctx, nil, f, postCondition, 1)
}

func increment(ctx context.Context, account eth.UserAccount, contract *bindings.Bethtest, val *big.Int) error {

	fmt.Printf("\n\x1b[37;1mIncrementing %v in the contract\x1b[0m\n", val.String())

	val.Add(val, big.NewInt(1))

	// Increment integer in the contract
	f := func(txOpts bind.TransactOpts) (*types.Transaction, error) {
		return contract.Increment(&txOpts)
	}

	// Post-condition: confirm that previous value has been incremented
	postCondition := func(ctx context.Context) bool {
		newVal, err := read(ctx, account.EthClient(), contract)
		if err != nil {
			return false
		}
		return newVal.Cmp(val) >= 0
	}

	return account.Transact(ctx, nil, f, postCondition, 2)
}

func elementExists(ctx context.Context, conn eth.Client, contract *bindings.Bethtest, val *big.Int) (exists bool) {
	exists = false
	_ = conn.Get(ctx, &bind.CallOpts{}, func(callOpts *bind.CallOpts) (err error) {
		_, exists, err = contract.Get(callOpts, val)
		return
	})
	return
}

func size(ctx context.Context, conn eth.Client, contract *bindings.Bethtest) (size *big.Int, err error) {
	size = big.NewInt(0)
	err = conn.Get(ctx, &bind.CallOpts{}, func(callOpts *bind.CallOpts) (err error) {
		size, err = contract.Size(callOpts)
		return
	})
	return
}

func read(ctx context.Context, conn eth.Client, contract *bindings.Bethtest) (newVal *big.Int, err error) {
	err = conn.Get(ctx, &bind.CallOpts{}, func(callOpts *bind.CallOpts) (err error) {
		newVal, err = contract.Read(callOpts)
		return
	})

	fmt.Printf("[info] Value in contract is %v\n", newVal.String())
	return
}
