package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/co-go"
	eth "github.com/republicprotocol/eth-go"
	"github.com/republicprotocol/republic-go/crypto"
)

// To run: go run transfer.go ../keystores/nightly-keystore.json 1.0 ./keystores/nightly-keystore.json xxxxxxxx 3a5e0b1158ca9ce861a80c3049d347a3f1825db0
// Where xxxxxxxx is the passphrase to unlock the keystore.

// Keystore paths for testing.
// trader := ../keystores/keystore-0.json
// genesis := ../keystores/nightly-keystore.json

// Addresses to test with.
// traderAddr := "3a5e0b1158ca9ce861a80c3049d347a3f1825db0"
// genesisAddr := "d5b5b26521665cb37623dca0e49c553b41dbf076"

func main() {

	// Check if all expected values were provided
	if len(os.Args) != 5 {
		fmt.Println("\nInvalid number of arguments!\x1b[0m " +
			"\n\nPlease enter a \x1b[37;1mamount\x1b[0m, " +
			"the \x1b[37;1mkeystore path\x1b[0m, a \x1b[37;1mpassphrase\x1b[0m " +
			"to unlock the keystore, and \x1b[37;1mthe ethereum address\x1b[0m " +
			"of the account to transfer value to.\n\n\x1b" +
			"[31;1m[Usage] go run transfer.go <amount> <path/to/keystore/file> " +
			"<passphrase> <ethereum-address-of-receiver>\x1b[0m")
		return
	}

	// Parse command-line arguments
	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Println("\n\x1b[31;1mPlease enter a valid value for amount\x1b[0m")
		return
	}
	keystorePath := os.Args[2]
	passphrase := os.Args[3]
	addr := os.Args[4]

	// Decrypt keystore
	ks := crypto.Keystore{}
	keyin, err := os.Open(keystorePath)

	if err != nil {
		log.Println(err)
		return
	}
	json, err := ioutil.ReadAll(keyin)
	if err != nil {
		log.Println(err)
		return
	}

	var privKey *ecdsa.PrivateKey

	if err := ks.DecryptFromJSON(json, passphrase); err != nil {
		key, err := keystore.DecryptKey(json, passphrase)
		if err != nil {
			log.Println(err)
			return
		}
		privKey = key.PrivateKey
	} else {
		privKey = ks.EcdsaKey.PrivateKey
	}

	// Get a user account to perform transactions
	account, err := eth.NewUserAccount("https://kovan.infura.io", privKey)
	if err != nil {
		log.Println(err)
		return
	}

	// Context with 5 minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	to := common.HexToAddress(addr)

	co.ParForAll([]int{1, 2}, func(i int) {
		// Transfer `x` amount of Eth to the specified address
		value, _ := big.NewFloat(amount * math.Pow10(18)).Int(nil)
		if err := account.Transfer(ctx, to, value, int64(i+1)); err != nil {
			log.Println(err)
			return
		}
	})

}
