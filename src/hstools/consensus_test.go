package hstools

import (
	"math/big"
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

	hsdir := c.Next(new(big.Int).SetBytes(desc[0])).Bytes()
	if ToHex(hsdir) != "274D66DC037FE344C58371B17C606988CBC37DFB" {
		t.Fatalf("wrong hsdir: %s", hsdir)
	}
}
