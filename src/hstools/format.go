package hstools

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"math"
	"math/big"
	"strings"
	"time"
)

type Hash [20]byte

// Distance efficiently calculates the difference (b - a) mod 2^160, or
// distance a -> b on a 20-byte ring and stores it in d. a and b unchanged.
func (a *Hash) Distance(b, d *Hash) {
	var carry bool
	for i := len(a) - 1; i >= 0; i-- {
		B := b[i]
		if carry {
			B--
		}
		d[i] = B - a[i]
		carry = B < a[i] || (carry && B == math.MaxUint8)
	}
}

func SHA1(data []byte) *Hash {
	h := Hash(sha1.Sum(data))
	return &h
}

func ToBase32(b []byte) string {
	return strings.ToLower(base32.StdEncoding.EncodeToString(b))
}

func FromBase32(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(strings.ToUpper(s))
}

func ToHex(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}

func FromHex(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

func FromBase64(s string) ([]byte, error) {
	if r := len(s) % 4; r != 0 {
		s += strings.Repeat("=", 4-r)
	}
	return base64.StdEncoding.DecodeString(s)
}

func HourToTime(h Hour) time.Time {
	return time.Unix(int64(h*3600), 0)
}

func HashesToIntSlice(keys []Hash) []*big.Int {
	ints := make([]*big.Int, len(keys))
	for i, k := range keys {
		ints[i] = new(big.Int).SetBytes(k[:])
	}
	return ints
}

func IntsToHashSlice(ints []*big.Int) []*Hash {
	hashes := make([]*Hash, len(ints))
	for n, i := range ints {
		var k Hash
		b := i.Bytes()
		copy(k[len(k)-len(b):], b)
		hashes[n] = &k
	}
	return hashes
}

func IntToHash(i *big.Int) Hash {
	var k Hash
	b := i.Bytes()
	copy(k[len(k)-len(b):], b)
	return k
}

func HashToInt(k Hash) *big.Int {
	return new(big.Int).SetBytes(k[:])
}
