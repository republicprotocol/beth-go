package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Conn contains the client and the contracts deployed to it
type Conn struct {
	Client *ethclient.Client
}

// Connect to an infura network (Supported networks: mainnet and kovan).
func Connect(network string) (Conn, error) {
	uri := ""

	switch network {
	case "mainnet":
		uri = "https://mainnet.infura.io"
	case "kovan":
		uri = "https://kovan.infura.io"
	default:
		return Conn{}, fmt.Errorf("cannot connect to %s: unsupported", network)
	}

	ethclient, err := ethclient.Dial(uri)
	if err != nil {
		return Conn{}, err
	}

	return Conn{
		Client: ethclient,
	}, nil
}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (conn *Conn) PatchedWaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := bind.WaitMined(ctx, conn.Client, tx)
	if err != nil {
		return nil, err
	}
	// if receipt.Status != types.ReceiptStatusSuccessful {
	// 	return receipt, errors.New("transaction reverted")
	// }
	return receipt, nil
}

// PatchedWaitDeployed waits for a contract deployment transaction and returns the on-chain
// contract address when it is mined. It stops waiting when ctx is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (conn *Conn) PatchedWaitDeployed(ctx context.Context, tx *types.Transaction) (common.Address, error) {
	return bind.WaitDeployed(ctx, conn.Client, tx)

}

// BalanceOf returns the ethereum balance of the addr passed.
func (conn *Conn) BalanceOf(addr common.Address) (*big.Int, error) {
	return conn.Client.BalanceAt(context.Background(), addr, nil)
}

// Get will perform
func (conn *Conn) Get(ctx context.Context, f func(opts *bind.CallOpts) (interface{}, error)) (interface{}, error) {
	sleepFor := time.Duration(1000)
	for {

		select {
		case <-ctx.Done():
			return nil, errors.New("context timed out")
		default:
		}

		if val, err := f(&bind.CallOpts{}); err == nil {
			return val, err
		}
		time.Sleep(sleepFor)

		// Increase delay for next round.
		sleepFor = time.Duration(float64(sleepFor) * 1.6)
	}
}
