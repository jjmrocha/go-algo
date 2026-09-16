// Package token provides compact, URL-safe identifiers that encode 128 bits
// as a fixed-length, lowercase base-36 string.
//
// A token is always 25 characters drawn from 0-9 and a-z, left-padded with
// zeros, so tokens sort lexicographically in the same order as the 128-bit
// values they encode.
package token

import (
	"crypto/rand"
	"encoding/binary"
	"math/bits"
	"uuid"
)

const (
	alphabet   = "0123456789abcdefghijklmnopqrstuvwxyz"
	base       = uint64(len(alphabet))
	tokenLen   = 25
	maxEncoded = "f5lxx1zz5pnorynqglhzmsp33"
)

var readRandom = func() (b [16]byte) {
	// crypto/rand.Read never returns an error since Go 1.24.
	_, _ = rand.Read(b[:])
	return b
}

// New returns a token encoding 128 bits from a cryptographically secure
// random source. No UUID version or variant bits are set.
func New() string {
	return encode(readRandom())
}

// FromUUID returns the token encoding the 128 bits of u. Any UUID version is
// accepted and its bits are used unchanged.
func FromUUID(u uuid.UUID) string {
	return encode(u)
}

// Valid reports whether s is a well-formed token: exactly 25 characters from
// 0-9 and a-z whose value fits in 128 bits. It does not report whether s was
// produced by this package.
func Valid(s string) bool {
	if len(s) != tokenLen {
		return false
	}

	for i := range len(s) {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'z') {
			return false
		}
	}

	return s <= maxEncoded
}

func encode(b [16]byte) string {
	hi := binary.BigEndian.Uint64(b[:8])
	lo := binary.BigEndian.Uint64(b[8:])

	var buf [tokenLen]byte

	for i := tokenLen - 1; i >= 0; i-- {
		var rem uint64
		hi, rem = hi/base, hi%base
		lo, rem = bits.Div64(rem, lo, base)
		buf[i] = alphabet[rem]
	}

	return string(buf[:])
}
