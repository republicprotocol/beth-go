package eth

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/eth-go/utils"
)

// Client will have a connection to an ethereum client (specified by the url)
type Client struct {
	ethClient *ethclient.Client
	url       string
}

// Connect to an infura network (Supported networks: mainnet and kovan).
func Connect(url string) (Client, error) {

	ethClient, err := ethclient.Dial(url)
	if err != nil {
		return Client{}, err
	}

	return Client{
		ethClient: ethClient,
		url:       url,
	}, nil
}

// WaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
func (client *Client) WaitMined(
	ctx context.Context,
	tx *types.Transaction,
) (*types.Receipt, error) {
	return bind.WaitMined(ctx, client.ethClient, tx)
}

// Get will perform a read-only transaction on the ethereum blockchain.
func (client *Client) Get(
	ctx context.Context,
	callOpts *bind.CallOpts,
	f func(*bind.CallOpts) (interface{}, error),
) (val interface{}, err error) {

	sleepDurationMs := time.Duration(1000)

	// Keep retrying until the read-only transaction succeeds or until context
	// times out
	for {
		select {
		case <-ctx.Done():
			if err == nil {
				return val, ctx.Err()
			}
			return
		default:
		}

		if val, err = f(callOpts); err == nil {
			return
		}

		// If transaction errors, wait for sometime before retrying
		select {
		case <-ctx.Done():
			if err == nil {
				return val, ctx.Err()
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
func (client *Client) BalanceOf(
	ctx context.Context,
	addr common.Address,
	callOpts *bind.CallOpts,
) (*big.Int, error) {

	val, err := client.Get(
		ctx,
		callOpts,
		func(*bind.CallOpts) (interface{}, error) {
			return client.ethClient.BalanceAt(ctx, addr, nil)
		},
	)
	if err != nil {
		return big.NewInt(0), err
	}
	return val.(*big.Int), nil
}

// GetBlockNumberByTxHash retrieves tx's block number using the tx hash.
func (client *Client) GetBlockNumberByTxHash(
	ctx context.Context,
	hash string,
) (*big.Int, error) {

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

		response, err := utils.SendRequest(ctx, client.url, jsonStr, data)
		if err != nil {
			continue
		}
		data = response.(JSONResponse)

		if data.Result == (Result{}) || data.Result.BlockNumber == "" {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(5 * time.Millisecond):
			}
			continue
		}
		break
	}

	return utils.HexToBigInt(data.Result.BlockNumber)
}

// GetCurrentBlockNumber will retrieve the current block that is confirmed by
// infura.
func (client *Client) GetCurrentBlockNumber(
	ctx context.Context,
) (*big.Int, error) {

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

		response, err := utils.SendRequest(ctx, client.url, jsonStr, data)
		if err != nil {
			continue
		}
		data = response.(JSONResponse)

		if data.Result == (Result{}) || data.Result.Number == "" {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(5 * time.Millisecond):
			}
			continue
		}
		break
	}

	return utils.HexToBigInt(data.Result.Number)
}
