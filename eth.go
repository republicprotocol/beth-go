package eth

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Account which is associated with a private key can read/write to the
// ethereum blockchain.
type Account interface {
	// Address returns the ethereum address of the account holder.
	Address() common.Address

	// Transfer sends the specified value of Eth to the given address.
	Transfer(ctx context.Context, to common.Address, value *big.Int) error

	// Transact performs a write operation on the Ethereum blockchain. It will
	// first conduct a preConditionCheck and if the check passes, it will
	// repeatedly execute the transaction followed by a postConditionCheck,
	// until the transaction passes and the postConditionCheck returns true.
	Transact(preConditionCheck func() bool, f func() (types.Transaction, error), postConditionCheck func() bool) error
}

// NewAccount returns an Account object associated with the given private key.
func NewAccount(network string, privateKey *ecdsa.PrivateKey) (Account, error) {
	return NewEthAccount(network, privateKey)
}

// NewConn returns a new connection to ethereum network. This connection can be
// used for performing read-only operations on smart contracts and for reading
// balance of ethereum addresses.
func NewConn(network string) (Conn, error) {
	return Connect(network)
}
