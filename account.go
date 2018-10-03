package eth

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ethAccount struct {
	conn         Conn
	transactOpts *bind.TransactOpts
}

func NewEthAccount(network string, privateKey *ecdsa.PrivateKey) (*ethAccount, error) {

	conn, err := Connect(network)
	if err != nil {
		return nil, nil
	}

	return &ethAccount{
		conn:         conn,
		transactOpts: bind.NewKeyedTransactor(privateKey),
	}, nil
}

func (account *ethAccount) Address() common.Address {
	return account.transactOpts.From
}

func (account *ethAccount) Transfer(ctx context.Context, to common.Address, value *big.Int) error {
	// Unimplemented
	return nil
}

func (account *ethAccount) Transact(preConditionCheck func() bool, f func() (types.Transaction, error), postConditionCheck func() bool) error {
	// Unimplemented
	return nil
}

func (account *ethAccount) SuggestGasPrice() (*big.Int, error) {
	return account.conn.Client.SuggestGasPrice(context.Background())
}

func (account *ethAccount) EstimateGas() {
	// return account.conn.Client.EstimateGas()
}
