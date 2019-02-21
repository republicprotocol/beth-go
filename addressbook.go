package beth

import (
	"github.com/ethereum/go-ethereum/common"
)

type AddressBook map[string]common.Address

var MainnetAddressBook = AddressBook{
	"DarknodeRegistry": common.HexToAddress("0x34bd421C7948Bc16f826Fd99f9B785929b121633"),
	"DGX":              common.HexToAddress("0x4f3AfEC4E5a3F2A6a1A411DEF7D7dFe50eE057bF"),
	"TUSD":             common.HexToAddress("0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E"),
	"REN":              common.HexToAddress("0x21C482f153D0317fe85C60bE1F7fa079019fcEbD"),
	"WBTC":             common.HexToAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"),
	"ZRX":              common.HexToAddress("0xE41d2489571d322189246DaFA5ebDe1F4699F498"),
	"OMG":              common.HexToAddress("0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"),
	"USDC":             common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
	"GUSD":             common.HexToAddress("0x056fd409e1d7a124bd7017459dfea2f387b6d5cd"),
	"DAI":              common.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
	"PAX":              common.HexToAddress("0x8e870d67f660d95d5be530380d0ec0bd388289e1"),
	"ETHSwapContract":  common.HexToAddress("0x4Bc1d23a8c00Ac87c57B6a32d5fb82aA5346950d"),
	"WBTCSwapContract": common.HexToAddress("0x15c10c51d86a51021d0683b8359fb20a8ba40b45"),
	"RENSwapContract":  common.HexToAddress("0xd633db90e6b017484ac08b711ed9f641c038141e"),
	"TUSDSwapContract": common.HexToAddress("0x0850ce7608313a957a52ef1536cdc0d87189d0e3"),
	"DGXSwapContract":  common.HexToAddress("0x42021bb8b52eae41a4bfed6af824f0a0d752d312"),
	"ZRXSwapContract":  common.HexToAddress("0x8a92e4f744460cb42e346636c3db48a940327107"),
	"OMGSwapContract":  common.HexToAddress("0x0f980ffa044bc28a0352a2282136bb61e2460ee4"),
}

var RopstenAddressBook = AddressBook{}

var KovanAddressBook = AddressBook{
	"RenExOrderbook":   common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"RenExSettlement":  common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"DarknodeRegistry": common.HexToAddress("0x75Fa8349fc9C7C640A4e9F1A1496fBB95D2Dc3d5"),
	"ETHSwapContract":  common.HexToAddress("0x6bBBDba39B529826Af6c070c4534b86FF6960B38"),
	"WBTCSwapContract": common.HexToAddress("0xb1e68BC63d5DF27e35b78Cef6d645Aa21586426e"),
	"RENSwapContract":  common.HexToAddress("0x63633aC3a001deeeb9200521DE7176302D4919B6"),
	"TUSDSwapContract": common.HexToAddress("0x3e2f1427787A52302d280f9e5F1416556EAbdB77"),
	"OMGSwapContract":  common.HexToAddress("0xA080114700A510e4d03DafD23f77B7E923D7676f"),
	"ZRXSwapContract":  common.HexToAddress("0xA3f3A0015B419eDF74568A9462279d25AD11e3B6"),
	"DGXSwapContract":  common.HexToAddress("0x4423F9E32904aAD9Fb7eaf3cA81cc439996C7D09"),
	"USDCSwapContract": common.HexToAddress("0xDD2E1bD1c6C87FF95d757C800D3c51852701b713"),
	"GUSDSwapContract": common.HexToAddress("0xF9b2d53E1FA1e3eeD0078E135c03691b1f12E3AC"),
	"DAISwapContract":  common.HexToAddress("0xa16FC6B7a48FDF135a13731Ba3c634DadC97FBc1"),
	"PAXSwapContract":  common.HexToAddress("0x694F850b9FCA9C9fDdBFB7BF16F4e348a4B29c01"),
	"WBTC":             common.HexToAddress("0xA1D3EEcb76285B4435550E4D963B8042A8bffbF0"),
	"REN":              common.HexToAddress("0x2CD647668494c1B15743AB283A0f980d90a87394"),
	"ZRX":              common.HexToAddress("0x6EB628dCeFA95802899aD3A9EE0C7650Ac63d543"),
	"OMG":              common.HexToAddress("0x66497ba75dD127b46316d806c077B06395918064"),
	"USDC":             common.HexToAddress("0x3f0a4aed397c66d7b7dde1d170321f87656b14cc"),
	"GUSD":             common.HexToAddress("0xA9CF366E9fb4F7959452d7a17A6F88ee2A20e9DA"),
	"DAI":              common.HexToAddress("0xc4375b7de8af5a38a93548eb8453a498222c4ff2"),
	"TUSD":             common.HexToAddress("0x525389752ffe6487d33EF53FBcD4E5D3AD7937a0"),
	"DGX":              common.HexToAddress("0x7d6D31326b12B6CBd7f054231D47CbcD16082b71"),
	"PAX":              common.HexToAddress("0x3584087444dabf2e0d29284766142ac5c3a9a2b7"),
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
