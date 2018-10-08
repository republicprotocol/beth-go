package typeutils

import (
	"errors"
	"math/big"
)

// ErrCannotConvertToBigInt is returned when string cannot be parsed into a big.Int format.
var ErrCannotConvertToBigInt = errors.New("cannot convert hex string to int: invalid format")

// HexToBigInt will convert a hex value in string format to the corresponding big.Int value. For example : "0xFD6CE" will return big.Int(1038030).
func HexToBigInt(hex string) (*big.Int, error) {
	bigInt := big.NewInt(0)
	bigIntStr := hex[2:]
	bigInt, ok := bigInt.SetString(bigIntStr, 16)
	if !ok {
		return bigInt, ErrCannotConvertToBigInt
	}
	return bigInt, nil
}
