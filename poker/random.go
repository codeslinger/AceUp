// Copyright (c) Toby DiPasquale. See accompanying LICENSE file for
// detailed licensing information.
package poker

import (
	"crypto/rand"
	"math/big"
)

// randInt generates a pseudo-random number in the range [0,max). It returns
// the generated number on success or -1 if a number could not be generated or
// max was less than or equal to 0.
func randInt(max int) int {
	if max <= 0 {
		return -1
	}
	m := big.NewInt(int64(max))
	r, e := rand.Int(rand.Reader, m)
	if e != nil {
		return -1
	}
	return int(r.Int64() % int64(max))
}
