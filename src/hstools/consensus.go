package hstools

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const consensusFilename = "consensuses-2006-01/02/2006-01-02-15-00-00-consensus"

// Hour is just a Unix timestamp divided by 3600, a unique index for an hour
type Hour int64

type Consensus struct {
	Time     Hour
	Filename string
	Error    error
	H        *Hashring
}

// ReadConsensuses reads consensus files from a folder structure like
// DIR/consensuses-2011-02/04/2011-02-04-02-00-00-consensus and sends them
// on the returned channel. From since to until included.
func ReadConsensuses(dir string, since, until Hour) chan Consensus {
	ch := make(chan Consensus)
	go func() {
		for h := since; h <= until; h++ {
			filename := time.Unix(int64(h*3600), 0).Format(consensusFilename)
			filename = filepath.Join(dir, filename)
			c := Consensus{
				Time:     h,
				Filename: filename,
			}

			hashring, err := ParseConsensus(filename)
			if err != nil {
				c.Error = err
			} else {
				c.H = hashring
			}

			ch <- c
		}
		close(ch)
	}()
	return ch
}

// ParseConsensus parses a consensus file and extracts the HSDir Hashring
func ParseConsensus(filename string) (*Hashring, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var fingerprint string
	var hsdirs []*big.Int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b := scanner.Bytes()

		if bytes.Equal(b[:2], []byte("r ")) {
			fingerprint = string(bytes.SplitN(b, []byte(" "), 4)[2])
			continue
		}

		if bytes.Equal(b[:2], []byte("s ")) && bytes.Contains(b, []byte("HSDir")) {
			f, err := FromBase64(fingerprint)
			if err != nil {
				return nil, fmt.Errorf("%v (%s)", err, fingerprint)
			}
			n := new(big.Int).SetBytes(f)
			hsdirs = append(hsdirs, n)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return NewHashring(hsdirs), nil
}
