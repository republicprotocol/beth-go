package beth_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/beth-go/test"
	"github.com/republicprotocol/co-go"
)

var _ = Describe("contracts", func() {

	newAccount := func(network, keystorePath string, passphrase string) (beth.Account, error) {
		// Open keystore file
		keyin, err := os.Open(keystorePath)
		if err != nil {
			return nil, err
		}
		data, err := ioutil.ReadAll(keyin)
		if err != nil {
			return nil, err
		}

		// Parse the json data to a key object
		key, err := keystore.DecryptKey(data, passphrase)
		if err != nil {
			return nil, err
		}

		return beth.NewAccount(fmt.Sprintf("https://%s.infura.io", network), key.PrivateKey)
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
		return test.NewBethtest(contractAddr, bind.ContractBackend(account.EthClient()))
	}

	elementExists := func(ctx context.Context, conn beth.Client, contract *test.Bethtest, val *big.Int) (exists bool) {
		exists = false
		_ = conn.Get(ctx, func() (err error) {
			_, exists, err = contract.Get(&bind.CallOpts{}, val)
			return
		})
		return
	}

	read := func(ctx context.Context, conn beth.Client, contract *test.Bethtest) (*big.Int, error) {
		newVal, err := contract.Read(&bind.CallOpts{})
		if err != nil {
			return nil, err
		}
		fmt.Printf("[info] Value in contract is %v\n", newVal.String())
		return newVal, nil
	}

	setInt := func(account beth.Account, contract *test.Bethtest, val *big.Int, waitBlocks int64) error {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		// Set integer in contract
		f := func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return contract.Set(txOpts, val)
		}

		// Post-condition: Confirm that the integer has the new value
		postCondition := func() bool {
			newVal, err := read(ctx, account.Client(), contract)
			if err != nil {
				return false
			}
			return newVal.Cmp(val) == 0
		}

		return account.Transact(ctx, nil, f, postCondition, waitBlocks)
	}

	increment := func(account beth.Account, contract *test.Bethtest, val *big.Int, waitBlocks int64) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Minute)
		defer cancel()

		val.Add(val, big.NewInt(1))

		// Increment integer in the contract
		f := func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
			return contract.Increment(txOpts)
		}

		// Post-condition: confirm that previous value has been incremented
		postCondition := func() bool {
			newVal, err := read(ctx, account.Client(), contract)
			if err != nil {
				return false
			}
			return newVal.Cmp(val) >= 0
		}

		return account.Transact(ctx, nil, f, postCondition, waitBlocks)
	}

	appendToList := func(values []*big.Int, contract *test.Bethtest, account beth.Account, waitBlocks int64) []error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration((len(values)*10)+int(waitBlocks))*time.Minute)
		defer cancel()

		errs := make([]error, len(values))

		co.ParForAll(values, func(i int) {
			defer GinkgoRecover()
			val := values[i]

			fmt.Printf("\n\x1b[37;1mAppending %v to list\x1b[0m", val.String())

			// Pre-condition: Does element already exist in the list?
			preCondition := func() bool {
				exists := elementExists(ctx, account.Client(), contract, val)
				if exists {
					fmt.Printf("\n[warning] Element %v exists in list!\n", val.String())
				}

				return !exists
			}

			// Append to list
			f := func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
				return contract.Append(txOpts, val)
			}

			// Post-condition: Has element been added to the list?
			postCondition := func() bool {
				return elementExists(ctx, account.Client(), contract, val)
			}

			// Execute transaction
			errs[i] = account.Transact(ctx, preCondition, f, postCondition, waitBlocks)
		})

		return errs
	}

	size := func(ctx context.Context, conn beth.Client, contract *test.Bethtest) (size *big.Int, err error) {
		size = big.NewInt(0)
		err = conn.Get(ctx, func() (err error) {
			size, err = contract.Size(&bind.CallOpts{})
			return
		})
		return
	}

	deleteFromList := func(values []*big.Int, contract *test.Bethtest, account beth.Account, waitBlocks int64) []error {
		// Context
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration((len(values)*10)+int(waitBlocks))*time.Minute)
		defer cancel()

		errs := make([]error, len(values))

		co.ParForAll(values, func(i int) {
			defer GinkgoRecover()
			val := values[i]
			fmt.Printf("\n\x1b[37;1mDeleting %v from list\x1b[0m", val.String())

			// Pre-condition: is list is empty or is element absent in list?
			preCondition := func() bool {

				// Read size of list
				size, err := size(ctx, account.Client(), contract)
				if err != nil || size.Cmp(big.NewInt(0)) <= 0 {
					fmt.Println("\n[warning] list is empty!")
					return false
				}

				// Check if element is present
				exists := elementExists(ctx, account.Client(), contract, val)
				if !exists {
					fmt.Printf("\n[warning] %v is not present in list!\n", val.String())
				}

				return exists
			}

			// Remove element
			f := func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
				return contract.Remove(txOpts, val)
			}

			// Post-condition: Has element been deleted successfully?
			postCondition := func() bool {
				return !elementExists(ctx, account.Client(), contract, val)
			}

			// Execute delete tx
			errs[i] = account.Transact(ctx, preCondition, f, postCondition, waitBlocks)
		})

		return errs
	}

	randomValues := func(n int64) []*big.Int {
		values := []*big.Int{}
		uniqueValues := make(map[int]struct{})
		i := 0
		for i < int(n) {
			// Randomly create a value and append it to list
			integerValue := rand.Intn(10000) + 1
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

	loadAddressBook := func(network string) beth.AddressBook {
		switch network {
		case "ropsten":
			return beth.DefaultAddressBook(3)
		case "kovan":
			return beth.DefaultAddressBook(42)
		default:
			return beth.AddressBook{}
		}
	}

	rand.Seed(time.Now().Unix())
	testedNetworks := []string{"ropsten", "kovan"}

	keystorePaths := []string{"test/keystore.ropsten.json", "test/keystore.kovan.json"}
	addresses := []string{"3a5e0b1158ca9ce861a80c3049d347a3f1825db0", "6b9b3e47c4c73db44f6a34064b21da8c62692a8c"}

	tableParallelism := []struct {
		n, waitBlocks int64
	}{
		{1, 3},
		// {2, 3},
		// {4, 2},
		// {8, 1},
		// {16, 1},
	}

	for _, network := range testedNetworks {
		network := network

		for ind, entry := range tableParallelism {
			n := entry.n
			waitBlocks := entry.waitBlocks

			// restrict ropsten tests
			if network == "ropsten" {
				if ind != 0 {
					continue
				}
				waitBlocks = 0
			}

			Context(fmt.Sprintf("when modifying an integer %v times in a contract deployed on %s", n, network), func() {

				It("should write to the contract and not return an error", func() {
					for i := 0; i < int(n); i++ {

						account, err := newAccount(network, fmt.Sprintf("test/keystore.%s.json", network), os.Getenv("passphrase"))
						Expect(err).ShouldNot(HaveOccurred())

						contract, err := bethTest(network, account)
						Expect(err).ShouldNot(HaveOccurred())

						// Generate random value
						val := big.NewInt(int64(rand.Intn(100)))

						fmt.Printf("\n\x1b[37;1mSetting integer %v in the contract on %s\n\x1b[0m", val.String(), network)

						nonceBefore, err := account.EthClient().NonceAt(context.Background(), account.Address(), nil)
						Expect(err).ShouldNot(HaveOccurred())

						// Set value in the contract
						err = setInt(account, contract, val, waitBlocks)
						Expect(err).ShouldNot(HaveOccurred())

						nonceMid, err := account.EthClient().NonceAt(context.Background(), account.Address(), nil)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(nonceMid - nonceBefore).Should(Equal(uint64(1)))

						fmt.Printf("\n\x1b[37;1mIncrementing %v in the contract on %s\x1b[0m\n", val.String(), network)

						// Increment the value in the contract
						err = increment(account, contract, val, waitBlocks)
						Expect(err).ShouldNot(HaveOccurred())

						nonceAfter, err := account.EthClient().NonceAt(context.Background(), account.Address(), nil)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(nonceAfter - nonceMid).Should(Equal(uint64(1)))
					}
				})
			})

			Context(fmt.Sprintf("when updating a list with %v elements in a contract deployed on %s", n, network), func() {

				It("should write to the contract and not return an error", func() {

					account, err := newAccount(network, fmt.Sprintf("test/keystore.%s.json", network), os.Getenv("passphrase"))
					Expect(err).ShouldNot(HaveOccurred())
					contract, err := bethTest(network, account)
					Expect(err).ShouldNot(HaveOccurred())

					// Retrieve original length of array
					originalLength, err := size(context.Background(), account.Client(), contract)
					Expect(err).ShouldNot(HaveOccurred())
					fmt.Printf("\n\x1b[37;1mThe size of array on %s is %v\n\x1b[0m", network, originalLength.String())

					// Append randomly generated values to a list maintained by a smart contract
					values := randomValues(n)
					errs := appendToList(values, contract, account, waitBlocks)
					err = handleErrors(errs)
					if err != nil && err != beth.ErrPreConditionCheckFailed {
						Expect(err).ShouldNot(HaveOccurred())
					}

					// Attempt to add a previously added item again
					errs = appendToList(values[:1], contract, account, waitBlocks)
					err = handleErrors(errs)
					Expect(err).Should(HaveOccurred())
					Expect(err).Should(Equal(beth.ErrPreConditionCheckFailed))

					// Attempt to delete all newly added elements from the list
					errs = deleteFromList(values, contract, account, waitBlocks)
					Expect(handleErrors(errs)).ShouldNot(HaveOccurred())

					// Attempt to delete a value that does not exist in the list
					errs = deleteFromList(values[:1], contract, account, waitBlocks)
					err = handleErrors(errs)
					Expect(err).Should(HaveOccurred())
					Expect(err).Should(Equal(beth.ErrPreConditionCheckFailed))

					// Retrieve length of array after deleting the newly added elements
					newLength, err := size(context.Background(), account.Client(), contract)
					Expect(err).ShouldNot(HaveOccurred())
					fmt.Printf("\n\x1b[37;1mThe new size of array on %s is %v\n\x1b[0m", network, newLength.String())

					// The new length must not be greater than the original length
					Expect(newLength.Cmp(originalLength)).To(BeNumerically("<=", 0))
				})
			})

			Context(fmt.Sprintf("when transferring eth from one account to an ethereum address on %s", network), func() {

				It("should successfully transfer eth and not return an error", func() {
					// Context with 10 minute timeout
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
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
						if err := account.Transfer(ctx, toAddrs[i], value, waitBlocks); err != nil {
							Expect(err).ShouldNot(HaveOccurred())
						}
					})
				})
			})

			Context(fmt.Sprintf("when retrieving addresses on %s", network), func() {

				It("should successfully return the address of RenExOrderbook", func() {
					addrBook := loadAddressBook(network)
					account, err := newAccount(network, keystorePaths[0], os.Getenv("passphrase"))
					Expect(err).ShouldNot(HaveOccurred())
					renExOrderbook, err := account.ReadAddress("RenExOrderbook")
					Expect(renExOrderbook.String()).Should(Equal(addrBook["RenExOrderbook"].String()))
				})

				It("should successfully return the address of RenExSettlement", func() {
					addrBook := loadAddressBook(network)
					account, err := newAccount(network, keystorePaths[0], os.Getenv("passphrase"))
					Expect(err).ShouldNot(HaveOccurred())
					renExSettlement, err := account.ReadAddress("RenExSettlement")
					Expect(renExSettlement.String()).Should(Equal(addrBook["RenExSettlement"].String()))
				})

				It("should successfully return the address of ERC20:WBTC", func() {
					addrBook := loadAddressBook(network)
					account, err := newAccount(network, keystorePaths[0], os.Getenv("passphrase"))
					Expect(err).ShouldNot(HaveOccurred())
					ERC20WBTC, err := account.ReadAddress("ERC20:WBTC")
					Expect(ERC20WBTC.String()).Should(Equal(addrBook["ERC20:WBTC"].String()))
				})

				It("should successfully return the address of Swapper:ETH", func() {
					addrBook := loadAddressBook(network)
					account, err := newAccount(network, keystorePaths[0], os.Getenv("passphrase"))
					Expect(err).ShouldNot(HaveOccurred())
					SwapperETH, err := account.ReadAddress("Swapper:ETH")
					Expect(SwapperETH.String()).Should(Equal(addrBook["Swapper:ETH"].String()))
				})

				It("should successfully return the address of Swapper:WBTC", func() {
					addrBook := loadAddressBook(network)
					account, err := newAccount(network, keystorePaths[0], os.Getenv("passphrase"))
					Expect(err).ShouldNot(HaveOccurred())
					SwapperWBTC, err := account.ReadAddress("Swapper:WBTC")
					Expect(SwapperWBTC.String()).Should(Equal(addrBook["Swapper:WBTC"].String()))
				})
			})

			Context("when signing messages", func() {
				It("should successfully sign a message", func() {
					account, err := newAccount(network, keystorePaths[0], os.Getenv("passphrase"))
					Expect(err).ShouldNot(HaveOccurred())

					msgHash := crypto.Keccak256([]byte("Message"))
					sig, err := account.Sign(msgHash)
					Expect(err).ShouldNot(HaveOccurred())

					publicKey, err := crypto.SigToPub(msgHash, sig)
					Expect(err).ShouldNot(HaveOccurred())

					signerAddress := crypto.PubkeyToAddress(*publicKey)
					Expect(signerAddress.String()).Should(Equal(account.Address().String()))
				})
			})
		}
	}
})
