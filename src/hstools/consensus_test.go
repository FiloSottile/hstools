package hstools

import (
	"encoding/base32"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestParseConsensus(t *testing.T) {
	c, err := ParseConsensus("../../test/2015-04-11-19-00-00-consensus")
	if err != nil {
		t.Fatal(err)
	}
	if c.Len() != 2983 {
		t.Fatalf("wrong number of points: %d", c.Len())
	}

	tt, _ := time.Parse(time.RFC3339, "2015-04-11T19:30:00Z")
	desc, err := OnionToDescID("facebookcorewwwi.onion", tt)

	descA, err := base32.StdEncoding.DecodeString(strings.ToUpper(desc[0]))
	if err != nil {
		t.Fatal(err)
	}
	hsdir := hex.EncodeToString(c.Next(new(big.Int).SetBytes(descA)).Bytes())
	if hsdir != "274d66dc037fe344c58371b17c606988cbc37dfb" {
		t.Fatalf("wrong hsdir: %s", hsdir)
	}
}
