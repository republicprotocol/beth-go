package beth_test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
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
	co "github.com/republicprotocol/co-go"
	"github.com/republicprotocol/republic-go/crypto"
)

var _ = Describe("contracts", func() {

	newAccount := func(network, keystorePath string, passphrase string) (beth.Account, error) {

		// Decrypt keystore
		ks := crypto.Keystore{}

		// Open keystore file
		keyin, err := os.Open(keystorePath)
		if err != nil {
			return nil, err
		}
		json, err := ioutil.ReadAll(keyin)
		if err != nil {
			return nil, err
		}

		var privKey *ecdsa.PrivateKey

		// Decrypt private key using keystore and passphrase
		if err := ks.DecryptFromJSON(json, passphrase); err != nil {
			key, err := keystore.DecryptKey(json, passphrase)
			if err != nil {
				return nil, err
			}
			privKey = key.PrivateKey
		} else {
			privKey = ks.EcdsaKey.PrivateKey
		}

		// Return a user account to perform transactions
		account, err := beth.NewAccount(fmt.Sprintf("https://%s.infura.io", network), privKey)
		if err != nil {
			return nil, err
		}

		return account, nil
	}

	bethTest := func(network string, account beth.Account) (*test.Bethtest, error) {
		// Ropsten : 0x46bcff69b2d5a677c40c05c1f034ef7bf0ee4742
		// Kovan : 0x055d30956deea82bfe0f99a2771dcc36a18dc9bb
		contractAddr := common.Address{}
		switch network {
		case "kovan":
			contractAddr = common.HexToAddress("0xaf866d7f173115e4cd5401f2abadb5b26eae8c32")
		case "ropsten":
			contractAddr = common.HexToAddress("0xd842576402d06f9985f407a5fd74e3eb06584110")
		default:
			return nil, errors.New("invalid infura network")
		}
		// Get contract
		return test.NewBethtest(contractAddr, bind.ContractBackend(account.EthClient().EthClient))
	}

	elementExists := func(ctx context.Context, conn beth.Client, contract *test.Bethtest, val *big.Int) (exists bool) {
		exists = false
		_ = conn.Get(ctx, &bind.CallOpts{}, func(callOpts *bind.CallOpts) (err error) {
			_, exists, err = contract.Get(callOpts, val)
			return
		})
		return
	}

	appendToList := func(values []*big.Int, contract *test.Bethtest, account beth.Account, numConfirms int64) []error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(len(values)+int(numConfirms)+3)*time.Minute)
		defer cancel()

		errs := make([]error, len(values))

		co.ParForAll(values, func(i int) {
			defer GinkgoRecover()
			val := values[i]

			fmt.Printf("\n\x1b[37;1mAppending %v to list\x1b[0m", val.String())

			// Pre-condition: Does element already exist in the list?
			preCondition := func(ctx context.Context) bool {
				exists := elementExists(ctx, account.EthClient(), contract, val)
				if exists {
					fmt.Printf("\n[warning] Element %v exists in list!\n", val.String())
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
			errs[i] = account.Transact(ctx, preCondition, f, postCondition, numConfirms)
		})

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

	deleteFromList := func(values []*big.Int, contract *test.Bethtest, account beth.Account, numConfirms int64) []error {
		// Context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(len(values)+int(numConfirms)+3)*time.Minute)
		defer cancel()

		errs := make([]error, len(values))

		co.ParForAll(values, func(i int) {
			defer GinkgoRecover()
			val := values[i]
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
			errs[i] = account.Transact(ctx, preCondition, f, postCondition, numConfirms)
		})
		return errs
	}

	randomValues := func(n int64) []*big.Int {
		rand.Seed(time.Now().Unix())

		values := []*big.Int{}
		uniqueValues := make(map[int]struct{})
		i := 0
		for i < int(n) {
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

	testedNetworks := []string{"kovan", "ropsten"}

	keystorePaths := []string{"test/keystore.ropsten.json", "test/keystore.kovan.json"}
	addresses := []string{"3a5e0b1158ca9ce861a80c3049d347a3f1825db0", "6b9b3e47c4c73db44f6a34064b21da8c62692a8c"}

	tableParallelism := []struct {
		n, numConfirms int64
	}{
		{1, 3},
		{2, 3},
		{4, 2},
		{8, 1},
		{16, 1},
	}

	for _, network := range testedNetworks {
		network := network

		Context(fmt.Sprintf("when updating a list in a contract deployed on %s", network), func() {

			for _, entry := range tableParallelism {
				n := entry.n
				numConfirms := entry.numConfirms

				It(fmt.Sprintf("should write %v elements to the contract and not return an error", n), func() {
					account, err := newAccount(network, fmt.Sprintf("test/keystore.%s.json", network), os.Getenv("passphrase"))
					Expect(err).ShouldNot(HaveOccurred())
					contract, err := bethTest(network, account)
					Expect(err).ShouldNot(HaveOccurred())

					// Retrieve original length of array
					originalLength, err := size(context.Background(), account.EthClient(), contract)
					Expect(err).ShouldNot(HaveOccurred())
					fmt.Printf("\n\x1b[37;1mThe size of array on %s is %v\n\x1b[0m", network, originalLength.String())

					// Append randomly generated values to a list maintained by a smart contract
					values := randomValues(n)
					errs := appendToList(values, contract, account, numConfirms)
					Expect(handleErrors(errs)).ShouldNot(HaveOccurred())

					// Attempt to add a previously added item again
					errs = appendToList(values[:1], contract, account, numConfirms)
					err = handleErrors(errs)
					Expect(err).Should(HaveOccurred())
					Expect(err).Should(Equal(beth.ErrorPreConditionCheckFailed))

					// Attempt to delete all newly added elements from the list
					errs = deleteFromList(values, contract, account, numConfirms)
					Expect(handleErrors(errs)).ShouldNot(HaveOccurred())

					// Attempt to delete a value that does not exist in the list
					errs = deleteFromList(values[:1], contract, account, numConfirms)
					err = handleErrors(errs)
					Expect(err).Should(HaveOccurred())
					Expect(err).Should(Equal(beth.ErrorPreConditionCheckFailed))

					// Retrieve length of array after deleting the newly added elements
					newLength, err := size(context.Background(), account.EthClient(), contract)
					Expect(err).ShouldNot(HaveOccurred())
					fmt.Printf("\n\x1b[37;1mThe new size of array on %s is %v\n\x1b[0m", network, newLength.String())

					// The new length must not be greater than the original length
					Expect(newLength.Cmp(originalLength)).To(BeNumerically("<=", 0))
				})
			}
		})

		Context(fmt.Sprintf("when transferring eth from one account to an ethereum address on %s", network), func() {

			It("should successfully transfer and not return an error", func() {
				// Context with 5 minute timeout
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				defer cancel()

				toAddrs := []common.Address{}
				for i := 0; i < 2; i++ {
					toAddrs = append(toAddrs, common.HexToAddress(addresses[i]))
				}

				co.ParForAll(toAddrs, func(i int) {

					account, err := newAccount(network, keystorePaths[i], os.Getenv("passphrase"))
					Expect(err).ShouldNot(HaveOccurred())
					// Transfer 1 Eth to the other account's address
					value, _ := big.NewFloat(1 * math.Pow10(18)).Int(nil)
					if err := account.Transfer(ctx, toAddrs[i], value, int64(i+1)); err != nil {
						Expect(err).ShouldNot(HaveOccurred())
					}
				})
			})
		})
	}
})
