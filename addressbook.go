package beth

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
)

type AddressBook map[string]common.Address

type DefaultAddressBooks struct {
	Mainnet []Address `json:"mainnet"`
	Ropsten []Address `json:"ropsten"`
	Kovan   []Address `json:"kovan"`
}

type Address struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func DefaultAddressBook(network int64) (AddressBook, error) {
	defaultAddrBooks := DefaultAddressBooks{}
	addrs := []Address{}
	addrBook := AddressBook{}

	defaultAddrBookData, err := ioutil.ReadFile("./addressbook.json")
	if err != nil {
		return addrBook, err
	}

	if err := json.Unmarshal(defaultAddrBookData, &defaultAddrBooks); err != nil {
		return addrBook, err
	}

	switch network {
	case 1:
		addrs = defaultAddrBooks.Mainnet
	case 3:
		addrs = defaultAddrBooks.Ropsten
	case 42:
		addrs = defaultAddrBooks.Kovan
	default:
		return addrBook, nil
	}

	for _, addr := range addrs {
		addrBook[addr.Name] = common.HexToAddress(addr.Address)
	}

	return addrBook, nil
}
