package hstools

import (
	"math/big"

	"git.torproject.org/user/phw/zoossh.git"
)

// ParseHashring parses a consensus file and extracts the HSDir Hashring
func ParseConsensus(filename string) (*Hashring, error) {
	c, err := zoossh.ParseConsensusFile(filename)
	if err != nil {
		return nil, err
	}

	var hsdirs []*big.Int
	for fingerprint, getStatus := range c.RouterStatuses {
		if !getStatus().Flags.HSDir {
			continue
		}
		n, ok := new(big.Int).SetString(fingerprint, 16)
		if !ok {
			panic("failed to parse hex")
		}
		hsdirs = append(hsdirs, n)
	}

	return NewHashring(hsdirs), nil
}
