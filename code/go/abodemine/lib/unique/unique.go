// Package unique implements functions that return random-ish strings
package unique

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// gen copied from https://github.com/nats-io/gnatsd/blob/master/util/mkpasswd.go
func gen(l int, ch []byte) []byte {
	b := make([]byte, l)
	max := big.NewInt(int64(len(ch)))

	for i := range b {
		ri, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(fmt.Sprintf("Error producing random integer: %v\n", err))
		}
		b[i] = ch[int(ri.Int64())]
	}

	return b
}
