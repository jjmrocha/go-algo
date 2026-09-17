// Package token provides compact, URL-safe identifiers that encode 128 bits
// as a fixed-length, lowercase base-36 string.
//
// A token is always 25 characters drawn from 0-9 and a-z, left-padded with
// zeros, so tokens sort lexicographically in the same order as the 128-bit
// values they encode.
package token

import (
	"encoding/binary"
	"errors"
	"math/bits"
	"uuid"
)

const (
	alphabet   = "0123456789abcdefghijklmnopqrstuvwxyz"
	base       = uint64(len(alphabet))
	tokenLen   = 25
	maxEncoded = "f5lxx1zz5pnorynqglhzmsp33"
)

// ErrInvalidToken is returned by [ToUUID] when its input is not a valid token.
var ErrInvalidToken = errors.New("invalid token")

// New returns a token encoding a new random UUID version 4, so [ToUUID]
// converts it back to a valid UUID v4.
func New() string {
	return FromUUID(uuid.NewV4())
}

// FromUUID returns the token encoding the 128 bits of u. Any UUID version is
// accepted and its bits are used unchanged.
func FromUUID(u uuid.UUID) string {
	return encode(u)
}

// ToUUID returns the UUID whose 128 bits s encodes. It is the inverse of
// [FromUUID]. If s is not a valid token (see [Valid]), it returns [uuid.Nil]
// and [ErrInvalidToken].
func ToUUID(s string) (uuid.UUID, error) {
	if !Valid(s) {
		return uuid.Nil(), ErrInvalidToken
	}

	return decode(s), nil
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

func decode(s string) (b [16]byte) {
	var hi, lo uint64

	for i := range len(s) {
		c := s[i]
		digit := uint64(c - '0')
		if c >= 'a' {
			digit = uint64(c-'a') + 10
		}

		var carry, add uint64
		carry, lo = bits.Mul64(lo, base)
		lo, add = bits.Add64(lo, digit, 0)
		hi = hi*base + carry + add
	}

	binary.BigEndian.PutUint64(b[:8], hi)
	binary.BigEndian.PutUint64(b[8:], lo)
	return b
}
