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
}

var RopstenAddressBook = AddressBook{}

var KovanAddressBook = AddressBook{
	"RenExOrderbook":   common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"RenExSettlement":  common.HexToAddress("0x0000000000000000000000000000000000000000"),
	"DarknodeRegistry": common.HexToAddress("0x75Fa8349fc9C7C640A4e9F1A1496fBB95D2Dc3d5"),
	"SwapperdETH":      common.HexToAddress("0x87E89520A5dab2d26d510c68A6bdAb93A252996F"),
	"SwapperdWBTC":     common.HexToAddress("0x3343180a98d580fc8E50a5493d7B5802829e4eff"),
	"SwapperdDGX":      common.HexToAddress("0xf86B4Ef43f7Bc1a1fb0CCe5F533A2267156A6f84"),
	"SwapperdTUSD":     common.HexToAddress("0x943E30f424b99eCCE86055409126b4C1be61C846"),
	"SwapperdREN":      common.HexToAddress("0x711c6dd63de71e5e0a1c5cfa09ad9079dc72e509"),
	"SwapperdZRX":      common.HexToAddress("0x6aD0aFbdA0011564adD307D3E7ef6d1669b2ECDA"),
	"SwapperdOMG":      common.HexToAddress("0x3aD7BFDDE63E56c6383A89635dA20E7aD5ff9801"),
	"WBTC":             common.HexToAddress("0xA1D3EEcb76285B4435550E4D963B8042A8bffbF0"),
	"DGX":              common.HexToAddress("0x932F4580B261e9781A6c3c102133C8fDd4503DFc"),
	"TUSD":             common.HexToAddress("0x525389752ffe6487d33EF53FBcD4E5D3AD7937a0"),
	"REN":              common.HexToAddress("0x2CD647668494c1B15743AB283A0f980d90a87394"),
	"ZRX":              common.HexToAddress("0x6EB628dCeFA95802899aD3A9EE0C7650Ac63d543"),
	"OMG":              common.HexToAddress("0x66497ba75dD127b46316d806c077B06395918064"),
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
