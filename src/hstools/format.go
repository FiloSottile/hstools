package hstools

import (
	"encoding/base32"
	"encoding/hex"
	"strings"
)

func ToBase32(b []byte) string {
	return strings.ToLower(base32.StdEncoding.EncodeToString(b))
}

func FromBase32(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(strings.ToUpper(s))
}

func ToHex(b []byte) string {
	return strings.ToUpper(hex.EncodeToString(b))
}
