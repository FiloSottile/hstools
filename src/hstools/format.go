package hstools

import (
	"encoding/base32"
	"encoding/base64"
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

func FromBase64(s string) ([]byte, error) {
	if r := len(s) % 4; r != 0 {
		s += strings.Repeat("=", 4-r)
	}
	return base64.StdEncoding.DecodeString(s)
}
