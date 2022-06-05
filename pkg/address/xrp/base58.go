package xrp

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
)

// Purloined from https://github.com/btcsuite/btcd/blob/master/btcutil/base58/base58.go

var (
	bigRadix10 = big.NewInt(58 * 58 * 58 * 58 * 58 * 58 * 58 * 58 * 58 * 58) // 58^10
	bigRadix   = [...]*big.Int{
		big.NewInt(0),
		big.NewInt(58),
		big.NewInt(58 * 58),
		big.NewInt(58 * 58 * 58),
		big.NewInt(58 * 58 * 58 * 58),
		big.NewInt(58 * 58 * 58 * 58 * 58),
		big.NewInt(58 * 58 * 58 * 58 * 58 * 58),
		big.NewInt(58 * 58 * 58 * 58 * 58 * 58 * 58),
		big.NewInt(58 * 58 * 58 * 58 * 58 * 58 * 58 * 58),
		big.NewInt(58 * 58 * 58 * 58 * 58 * 58 * 58 * 58 * 58),
		bigRadix10,
	}
)

// Base58Decode decodes a modified base58 string to a byte slice and checks checksum.
func Base58Decode(b, alphabet string) ([]byte, error) {
	if len(b) < 5 {
		return nil, fmt.Errorf("Base58 string too short: %s", b)
	}
	answer := big.NewInt(0)
	scratch := new(big.Int)

	// Calculating with big.Int is slow for each iteration.
	//    x += b58[b[i]] * j
	//    j *= 58
	//
	// Instead we can try to do as much calculations on int64.
	// We can represent a 10 digit base58 number using an int64.
	//
	// Hence we'll try to convert 10, base58 digits at a time.
	// The rough idea is to calculate `t`, such that:
	//
	//   t := b58[b[i+9]] * 58^9 ... + b58[b[i+1]] * 58^1 + b58[b[i]] * 58^0
	//   x *= 58^10
	//   x += t
	//
	// Of course, in addition, we'll need to handle boundary condition when `b` is not multiple of 58^10.
	// In that case we'll use the bigRadix[n] lookup for the appropriate power.
	for t := b; len(t) > 0; {
		n := len(t)
		if n > 10 {
			n = 10
		}

		total := uint64(0)
		for _, v := range t[:n] {
			tmp := strings.Index(alphabet, string(v))
			if tmp == -1 {
				return nil, fmt.Errorf("Bad Base58 string: %s", b)
			}
			total = total*58 + uint64(tmp)
		}

		answer.Mul(answer, bigRadix[n])
		scratch.SetUint64(total)
		answer.Add(answer, scratch)

		t = t[n:]
	}

	tmpval := answer.Bytes()

	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != alphabet[0] {
			break
		}
	}
	flen := numZeros + len(tmpval)
	val := make([]byte, flen)
	copy(val[numZeros:], tmpval)

	// Check checksum
	checksum := DoubleSha256(val[0 : len(val)-4])
	expected := val[len(val)-4:]
	if !bytes.Equal(checksum[0:4], expected) {
		return nil, fmt.Errorf("Bad Base58 checksum: %v expected %v", checksum, expected)
	}
	return val, nil
}

// Base58Encode encodes a byte slice to a modified base58 string.
func Base58Encode(b []byte, alphabet string) string {
	checksum := DoubleSha256(b)
	b = append(b, checksum[0:4]...)
	x := new(big.Int)
	x.SetBytes(b)

	// maximum length of output is log58(2^(8*len(b))) == len(b) * 8 / log(58)
	maxlen := int(float64(len(b))*1.365658237309761) + 1
	answer := make([]byte, 0, maxlen)
	mod := new(big.Int)
	for x.Sign() > 0 {
		// Calculating with big.Int is slow for each iteration.
		//    x, mod = x / 58, x % 58
		//
		// Instead we can try to do as much calculations on int64.
		//    x, mod = x / 58^10, x % 58^10
		//
		// Which will give us mod, which is 10 digit base58 number.
		// We'll loop that 10 times to convert to the answer.

		x.DivMod(x, bigRadix10, mod)
		if x.Sign() == 0 {
			// When x = 0, we need to ensure we don't add any extra zeros.
			m := mod.Int64()
			for m > 0 {
				answer = append(answer, alphabet[m%58])
				m /= 58
			}
		} else {
			m := mod.Int64()
			for i := 0; i < 10; i++ {
				answer = append(answer, alphabet[m%58])
				m /= 58
			}
		}
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, alphabet[0])
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}
	return string(answer)
}
