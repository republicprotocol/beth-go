package beth_test

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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	beth "github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/beth-go/test"
	"github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("contracts", func() {

	newAccount := func(network, keystorePath string, passphrase string) (beth.Account, *test.Bethtest, error) {
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
		account, err := beth.NewAccount(fmt.Sprintf("https://%s.infura.io", network), privKey)
		if err != nil {
			return nil, nil, err
		}

		// Get contract
		contract, err := test.NewBethtest(contractAddr, bind.ContractBackend(account.EthClient().EthClient))
		if err != nil {
			return nil, nil, err
		}

		return account, contract, nil
	}

	elementExists := func(ctx context.Context, conn beth.Client, contract *test.Bethtest, val *big.Int) (exists bool) {
		exists = false
		_ = conn.Get(ctx, &bind.CallOpts{}, func(callOpts *bind.CallOpts) (err error) {
			_, exists, err = contract.Get(callOpts, val)
			return
		})
		return
	}

	appendToList := func(values []*big.Int, contract *test.Bethtest, account beth.Account) []error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(len(values)+2)*time.Minute)
		defer cancel()

		errs := make([]error, len(values))

		for i, val := range values {

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
			errs[i] = account.Transact(ctx, preCondition, f, postCondition, 1)
		}

		return errs
	}

	size := func(ctx context.Context, conn beth.Client, contract *test.Bethtest) (size *big.Int, err error) {
		size = big.NewInt(0)
		err = conn.Get(ctx, &bind.CallOpts{}, func(callOpts *bind.CallOpts) (err error) {
			size, err = contract.Size(callOpts)
			return
		})
		return
	}

	deleteFromList := func(values []*big.Int, contract *test.Bethtest, account beth.Account) []error {
		// Context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(len(values)+2)*time.Minute)
		defer cancel()

		errs := make([]error, len(values))

		for i, val := range values {
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
			errs[i] = account.Transact(ctx, preCondition, f, postCondition, 1)
		}

		return errs
	}

	randomValues := func(n int) []*big.Int {
		rand.Seed(time.Now().Unix())

		values := []*big.Int{}
		uniqueValues := make(map[int]struct{})
		i := 0
		for i < n {
			// Randomly create a value and append it to list
			integerValue := rand.Intn(100) + 1
			if _, ok := uniqueValues[integerValue]; ok {
				continue
			}
			uniqueValues[integerValue] = struct{}{}

			val := big.NewInt(int64(integerValue))
			values = append(values, val)
			i++
		}
		return values
	}

	// handleErrors loop through a list of errors. Calculate how many errors is not
	// nil and return the last non-nil error if exists.
	handleErrors := func(errs []error) error {
		var counter int
		var err error

		for i := range errs {
			if errs[i] != nil {
				counter++
				err = errs[i]
			}
		}
		return err
	}

	Context("when modifying a list in the contract", func() {

		It("should update the contract and not return an error", func() {
			n := 1
			account, contract, err := newAccount("kovan", "test/keystore.ropsten.json", os.Getenv("passphrase"))
			Expect(err).ShouldNot(HaveOccurred())

			// Retrieve original length of array
			originalLength, err := size(context.Background(), account.EthClient(), contract)
			Expect(err).ShouldNot(HaveOccurred())

			// Append randomly generated values to a list maintained by a smart contract
			values := randomValues(n)
			errs := appendToList(values, contract, account)
			Expect(handleErrors(errs)).ShouldNot(HaveOccurred())

			// Attempt to add a previously added item again
			errs = appendToList(values[:1], contract, account)
			err = handleErrors(errs)
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(Equal(beth.ErrorPreConditionCheckFailed))

			// Attempt to delete all newly added elements from the list
			errs = deleteFromList(values, contract, account)
			Expect(handleErrors(errs)).ShouldNot(HaveOccurred())

			// Attempt to delete a value that does not exist in the list
			errs = deleteFromList(values[:1], contract, account)
			err = handleErrors(errs)
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(Equal(beth.ErrorPreConditionCheckFailed))

			// Retrieve length of array after deleting the newly added elements
			newLength, err := size(context.Background(), account.EthClient(), contract)
			Expect(err).ShouldNot(HaveOccurred())

			// The new length must not be greater than the original length
			Expect(newLength.Cmp(originalLength)).To(BeNumerically("<=", 0))
		})
	})
})
