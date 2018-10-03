package eth

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// UserAccount which is associated with a private key can read/write to the
// ethereum blockchain.
type UserAccount interface {
	// Address returns the ethereum address of the account holder.
	Address() common.Address

	// Transfer sends the specified value of Eth to the given address.
	Transfer(ctx context.Context, to common.Address, value *big.Int, confirmBlocks int64) error

	// Transact performs a write operation on the Ethereum blockchain. It will
	// first conduct a preConditionCheck and if the check passes, it will
	// repeatedly execute the transaction followed by a postConditionCheck,
	// until the transaction passes and the postConditionCheck returns true.
	Transact(ctx context.Context, preConditionCheck func() bool, f func(bind.TransactOpts) (*types.Transaction, error), postConditionCheck func(ctx context.Context) bool, confirmBlocks int64) error
}

// NewUserAccount returns an Account object associated with the given private key.
func NewUserAccount(network string, privateKey *ecdsa.PrivateKey) (UserAccount, error) {
	return NewAccount(network, privateKey)
}

// NewClient returns a new connection to ethereum network. This connection can be
// used for performing read-only operations on smart contracts and for reading
// balance of ethereum addresses.
func NewClient(network string) (Client, error) {
	return Connect(network)
}
