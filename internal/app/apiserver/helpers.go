package apiserver

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// big Hex to human-readable string with float. Used to Convert transaction.Value to human-readable float string
func BigHexToStr() func(string) string {
	conversionFactor := new(big.Float).SetUint64(1000000000000000000)

	return func(hexValue string) string {
		vTogether := new(big.Float).SetInt(hexutil.MustDecodeBig(hexValue))
		bigFloat := vTogether.Quo(vTogether, conversionFactor)
		return strings.TrimRight(fmt.Sprintf("%3.18f", bigFloat), "0")
	}
}
