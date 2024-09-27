package types

import (
	"math/big"

	"github.com/holiman/uint256"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// BytesToUint256 strictly returns nil if byte array is nil
func BytesToUint256(b []byte) *uint256.Int {
	if b == nil {
		return nil
	}
	return new(uint256.Int).SetBytes(b)
}

// BytesToBigInt strictly returns nil if byte array is nil
func BytesToBigInt(b []byte) *big.Int {
	if b == nil {
		return nil
	}
	return new(big.Int).SetBytes(b)
}

// BigIntToUint256 strictly returns nil if big int is nil
func BigIntToUint256(i *big.Int) *uint256.Int {
	if i == nil {
		return nil
	}
	return uint256.MustFromBig(i)
}

// BigIntToWrapperspbBytes strictly returns nil if big int is nil
func BigIntToWrapperspbBytes(i *big.Int) *wrapperspb.BytesValue {
	if i == nil {
		return nil
	}
	return wrapperspb.Bytes(i.Bytes())
}
