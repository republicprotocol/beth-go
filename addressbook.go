package beth

import (
	"github.com/ethereum/go-ethereum/common"
)

type AddressBook map[string]common.Address

type Address struct {
	Name    string
	Address string
}

var MainnetAddresses = []Address{}

var RopstenAddresses = []Address{
	Address{
		Name:    "RenExOrderbook",
		Address: "0xA1D3EEcb76285B4435550E4D963B8042A8b11111",
	},
	Address{
		Name:    "RenExSettlement",
		Address: "0xA1D3EEcb76285B4435550E4D963B8042A8b11111",
	},
	Address{
		Name:    "WBTC",
		Address: "0xA1D3EEcb76285B4435550E4D963B8042A8b11111",
	},
	Address{
		Name:    "SwapperdETH",
		Address: "0x2218fa20c33765e7e01671ee6aaca75fbaf3a974",
	},
	Address{
		Name:    "SwapperdWBTC",
		Address: "0x2218fa20c33765e7e01671ee6aaca75fbaf3a974",
	},
}

var KovanAddresses = []Address{
	Address{
		Name:    "RenExOrderbook",
		Address: "0xA1D3EEcb76285B4435550E4D963B8042A8b11111",
	},
	Address{
		Name:    "RenExSettlement",
		Address: "0xA1D3EEcb76285B4435550E4D963B8042A8b11111",
	},
	Address{
		Name:    "WBTC",
		Address: "0xA1D3EEcb76285B4435550E4D963B8042A8bffbF0",
	},
	Address{
		Name:    "SwapperdETH",
		Address: "0x2218fa20c33765e7e01671ee6aaca75fbaf3a974",
	},
	Address{
		Name:    "SwapperdWBTC",
		Address: "0x2218fa20c33765e7e01671ee6aaca75fbaf3a974",
	},
}

func DefaultAddressBook(network int64) AddressBook {
	addrs := []Address{}
	addrBook := AddressBook{}

	switch network {
	case 1:
		addrs = MainnetAddresses
	case 3:
		addrs = RopstenAddresses
	case 42:
		addrs = KovanAddresses
	default:
		return addrBook
	}

	for _, addr := range addrs {
		addrBook[addr.Name] = common.HexToAddress(addr.Address)
	}
	return addrBook
}
