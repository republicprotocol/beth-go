package beth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type erc20 struct {
	account *account
	cerc20  *CompatibleERC20
}

type ERC20 interface {
	BalanceOf(ctx context.Context, who common.Address) (*big.Int, error)
	Allowance(ctx context.Context, owner, spender common.Address) (*big.Int, error)
	Transfer(ctx context.Context, to common.Address, amount *big.Int) (string, error)
	Approve(ctx context.Context, spender common.Address, amount *big.Int) (string, error)
	TransferFrom(ctx context.Context, from, to common.Address, amount *big.Int) (string, error)
}

func (account *account) NewERC20(addressOrAlias string) (ERC20, error) {
	address, ok := account.addressBook[addressOrAlias]
	if !ok {
		address = common.HexToAddress(addressOrAlias)
	}
	compatibleERC20, err := NewCompatibleERC20(address, bind.ContractBackend(account.EthClient()))
	if err != nil {
		return nil, err
	}
	return &erc20{
		account: account,
		cerc20:  compatibleERC20,
	}, nil
}

func (erc20 *erc20) BalanceOf(ctx context.Context, who common.Address) (*big.Int, error) {
	client := erc20.account.Client()
	var balance *big.Int
	return balance, client.Get(ctx, func() error {
		bal, err := erc20.cerc20.BalanceOf(&bind.CallOpts{}, who)
		if err != nil {
			return err
		}
		balance = bal
		return nil
	})
}

func (erc20 *erc20) Allowance(ctx context.Context, owner, spender common.Address) (*big.Int, error) {
	client := erc20.account.Client()
	var allowance *big.Int
	return allowance, client.Get(ctx, func() error {
		alw, err := erc20.cerc20.Allowance(&bind.CallOpts{}, owner, spender)
		if err != nil {
			return err
		}
		allowance = alw
		return nil
	})
}

func (erc20 *erc20) Transfer(ctx context.Context, to common.Address, amount *big.Int) (string, error) {
	if amount == nil {
		balance, err := erc20.BalanceOf(ctx, erc20.account.Address())
		if err != nil {
			return "", err
		}
		amount = balance
	}
	var txHash string
	return txHash, erc20.account.Transact(
		ctx,
		nil,
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := erc20.cerc20.Transfer(tops, to, amount)
			if err != nil {
				return tx, err
			}
			txHash = tx.Hash().String()
			return tx, nil
		},
		nil,
		1,
	)
}

func (erc20 *erc20) Approve(ctx context.Context, spender common.Address, amount *big.Int) (string, error) {
	var txHash string
	return txHash, erc20.account.Transact(
		ctx,
		nil,
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := erc20.cerc20.Approve(tops, spender, amount)
			if err != nil {
				return tx, err
			}
			txHash = tx.Hash().String()
			return tx, nil
		},
		nil,
		1,
	)
}

func (erc20 *erc20) TransferFrom(ctx context.Context, from, to common.Address, amount *big.Int) (string, error) {
	var txHash string
	return txHash, erc20.account.Transact(
		ctx,
		nil,
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := erc20.cerc20.TransferFrom(tops, from, to, amount)
			if err != nil {
				return tx, err
			}
			txHash = tx.Hash().String()
			return tx, nil
		},
		nil,
		1,
	)
}
