package beth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ErrCannotConvertToBigInt is returned when string cannot be parsed into a
// big.Int format.
var ErrCannotConvertToBigInt = errors.New("cannot convert hex string to int: invalid format")

// Client will have a connection to an ethereum client (specified by the url)
type Client struct {
	ethClient *ethclient.Client
	addrBook  AddressBook
	url       string
}

// Connect to an infura network (Supported networks: mainnet and kovan).
func Connect(url string) (Client, error) {

	ethClient, err := ethclient.Dial(url)
	if err != nil {
		return Client{}, err
	}

	netID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return Client{}, err
	}

	return Client{
		ethClient: ethClient,
		addrBook:  DefaultAddressBook(netID.Int64()),
		url:       url,
	}, nil
}

// WriteAddress to the address book, overwrite if already exists
func (client *Client) WriteAddress(key string, address common.Address) {
	client.addrBook[key] = address
}

// ReadAddress from the address book, return an error if the address does not
// exist
func (client *Client) ReadAddress(key string) (common.Address, error) {
	if address, ok := client.addrBook[key]; ok {
		return address, nil
	}
	return common.Address{}, ErrAddressNotFound
}

// WaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
func (client *Client) WaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	return bind.WaitMined(ctx, client.ethClient, tx)
}

// Get will perform a read-only transaction on the ethereum blockchain.
func (client *Client) Get(ctx context.Context, f func() error) (err error) {

	sleepDurationMs := time.Duration(1000)

	// Keep retrying until the read-only transaction succeeds or until context
	// times out
	for {
		select {
		case <-ctx.Done():
			if err == nil {
				return ctx.Err()
			}
			return
		default:
		}

		if err = f(); err == nil {
			return
		}

		// If transaction errors, wait for sometime before retrying
		select {
		case <-ctx.Done():
			if err == nil {
				return ctx.Err()
			}
			return
		case <-time.After(sleepDurationMs * time.Millisecond):
		}

		// Increase delay for next round but saturate at 30s
		sleepDurationMs = time.Duration(float64(sleepDurationMs) * 1.6)
		if sleepDurationMs > 30000 {
			sleepDurationMs = 30000
		}
	}
}

// BalanceOf returns the ethereum balance of the addr passed.
func (client *Client) BalanceOf(ctx context.Context, addr common.Address) (val *big.Int, err error) {
	err = client.Get(ctx, func() (err error) {
		val, err = client.ethClient.BalanceAt(ctx, addr, nil)
		return
	})
	return
}

// EthClient returns the ethereum client connection.
func (client *Client) EthClient() *ethclient.Client {
	return client.ethClient
}

// TxBlockNumber retrieves tx's block number using the tx hash.
func (client *Client) TxBlockNumber(ctx context.Context, hash string) (*big.Int, error) {

	type Result struct {
		BlockNumber string `json:"blockNumber,omitempty"`
	}
	type JSONResponse struct {
		Result Result `json:"result,omitempty"`
	}
	var data JSONResponse

	var jsonStr = `{"jsonrpc":"2.0","method":"eth_getTransactionByHash",` +
		`"params":["` + hash + `"],"id":1}`

	// Keep retrying until a block number is returned or until context times out
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		response, err := sendInfuraRequest(ctx, client.url, jsonStr)
		if err != nil {
			continue
		}
		err = json.Unmarshal(response, &data)

		if err != nil || data.Result == (Result{}) || data.Result.BlockNumber == "" {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(5 * time.Millisecond):
			}
			continue
		}
		break
	}

	return hexToBigInt(data.Result.BlockNumber)
}

// CurrentBlockNumber will retrieve the current block that is confirmed by
// infura.
func (client *Client) CurrentBlockNumber(ctx context.Context) (*big.Int, error) {

	type Result struct {
		Number string `json:"number,omitempty"`
	}
	type JSONResponse struct {
		Result Result `json:"result,omitempty"`
	}
	var data JSONResponse

	var jsonStr = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber",` +
		`"params":["latest", false],"id":1}`

	// Keep retrying until a block number is returned or until context times out
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		response, err := sendInfuraRequest(ctx, client.url, jsonStr)
		if err != nil {
			continue
		}
		err = json.Unmarshal(response, &data)

		if err != nil || data.Result == (Result{}) || data.Result.Number == "" {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(5 * time.Millisecond):
			}
			continue
		}
		break
	}

	return hexToBigInt(data.Result.Number)
}

// hexToBigInt will convert a hex value in string format to the corresponding
// big.Int value. For example : "0xFD6CE" will return big.Int(1038030).
func hexToBigInt(hex string) (*big.Int, error) {
	bigInt := big.NewInt(0)
	bigIntStr := hex[2:]
	bigInt, ok := bigInt.SetString(bigIntStr, 16)
	if !ok {
		return bigInt, ErrCannotConvertToBigInt
	}
	return bigInt, nil
}

// sendInfuraRequest will send a request to infura and return the unmarshalled data
// back to the caller. It will retry until a valid response is returned, or
// until the context times out.
func sendInfuraRequest(ctx context.Context, url string, request string) (body []byte, err error) {

	sleepDurationMs := time.Duration(1000)

	// Retry until a valid response is returned or until context times out
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if body, err = func() ([]byte, error) {
			// Create a new http POST request
			req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(request)))
			if err != nil {
				return nil, err
			}

			// Send http POST request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}

			// Decode response body
			return func() ([]byte, error) {
				defer resp.Body.Close()

				// Check status
				if resp.StatusCode != http.StatusOK {
					return nil, fmt.Errorf("unexpected status %v", resp.StatusCode)
				}
				// Check body
				if resp.Body != nil {
					return ioutil.ReadAll(resp.Body)
				}
				return nil, fmt.Errorf("response body is nil")
			}()
		}(); err == nil {
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(sleepDurationMs * time.Millisecond):

		}

		// Increase delay for next round but saturate at 30s
		sleepDurationMs = time.Duration(float64(sleepDurationMs) * 1.6)
		if sleepDurationMs > 30000 {
			sleepDurationMs = 30000
		}
	}
	return
}
