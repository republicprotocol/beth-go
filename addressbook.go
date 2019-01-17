package beth

import (
	"github.com/ethereum/go-ethereum/common"
)

type AddressBook map[string]common.Address

var MainnetAddressBook = AddressBook{
	"DGX":  common.HexToAddress("0x4f3AfEC4E5a3F2A6a1A411DEF7D7dFe50eE057bF"),
	"TUSD": common.HexToAddress("0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E"),
	"REN":  common.HexToAddress("0x21C482f153D0317fe85C60bE1F7fa079019fcEbD"),
	"WBTC": common.HexToAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"),
	"ZRX":  common.HexToAddress("0xE41d2489571d322189246DaFA5ebDe1F4699F498"),
	"OMG":  common.HexToAddress("0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"),
	"USDC": common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
	"GUSD": common.HexToAddress("0x056fd409e1d7a124bd7017459dfea2f387b6d5cd"),
	"DAI":  common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
}

var RopstenAddressBook = AddressBook{}

var KovanAddressBook = AddressBook{
	"RenExOrderbook":   common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"RenExSettlement":  common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"DarknodeRegistry": common.HexToAddress("0x75Fa8349fc9C7C640A4e9F1A1496fBB95D2Dc3d5"),
	"ETHSwapContract":  common.HexToAddress("0x94ab22cffb9cc1ee4097ff17ef9c02fbb26fdfa4"),
	"WBTCSwapContract": common.HexToAddress("0xad29a79ae4863ea605b14a6be730e29bcd5f2294"),
	"RENSwapContract":  common.HexToAddress("0x7708a58e7c1fdc6d9e092e4270f15f30ffffbbaf"),
	"ZRXSwapContract":  common.HexToAddress("0x3f9032cd9bb2233694e7b51ada014345216c6a90"),
	"OMGSwapContract":  common.HexToAddress("0x07ca0635574a191a7e63ce1d67c295cafbcf1e87"),
	"USDCSwapContract": common.HexToAddress("0x4cc61223b5308ff6b48b43d0a2c425b5f11f9bca"),
	"GUSDSwapContract": common.HexToAddress("0xb9efa9dc4306ae3a5abe652db65361d13765c0d6"),
	"DAISwapContract":  common.HexToAddress("0x233fdd9253fda1f616bfd382633ca74e4824b703"),
	"TUSDSwapContract": common.HexToAddress("0x61ba8c2d07d701056df3e6038d9abc25f6b86da6"),
	"DGXSwapContract":  common.HexToAddress("0x84715415755d44e90f78639e27379017e8a038ef"),
	"WBTC":             common.HexToAddress("0xA1D3EEcb76285B4435550E4D963B8042A8bffbF0"),
	"REN":              common.HexToAddress("0x2CD647668494c1B15743AB283A0f980d90a87394"),
	"ZRX":              common.HexToAddress("0x6EB628dCeFA95802899aD3A9EE0C7650Ac63d543"),
	"OMG":              common.HexToAddress("0x66497ba75dD127b46316d806c077B06395918064"),
	"USDC":             common.HexToAddress("0x3f0a4aed397c66d7b7dde1d170321f87656b14cc"),
	"GUSD":             common.HexToAddress("0xA9CF366E9fb4F7959452d7a17A6F88ee2A20e9DA"),
	"DAI":              common.HexToAddress("0xc4375b7de8af5a38a93548eb8453a498222c4ff2"),
	"TUSD":             common.HexToAddress("0x525389752ffe6487d33EF53FBcD4E5D3AD7937a0"),
	"DGX":              common.HexToAddress("0x932F4580B261e9781A6c3c102133C8fDd4503DFc"),
}

func DefaultAddressBook(network int64) AddressBook {
	switch network {
	case 1:
		return MainnetAddressBook
	case 3:
		return RopstenAddressBook
	case 42:
		return KovanAddressBook
	default:
		return AddressBook{}
	}
}
