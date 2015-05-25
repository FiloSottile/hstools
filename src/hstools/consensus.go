package hstools

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const consensusFilename = "consensuses-2006-01/02/2006-01-02-15-00-00-consensus"

// Hour is just a Unix timestamp divided by 3600, a unique index for an hour
type Hour int32

type Consensus struct {
	Time     Hour
	Filename string
	Error    error
	K        []Hash
}

// ReadConsensuses reads consensus files from a folder structure like
// DIR/consensuses-2011-02/04/2011-02-04-02-00-00-consensus and sends them
// on the returned channel. From since to until included.
func ReadConsensuses(dir string, since, until Hour) chan *Consensus {
	ch := make(chan *Consensus)
	go func() {
		for h := since; h <= until; h++ {
			filename := HourToTime(h).Format(consensusFilename)
			filename = filepath.Join(dir, filename)
			c := &Consensus{
				Time:     h,
				Filename: filename,
			}

			keys, err := ParseConsensus(filename)
			if err != nil {
				c.Error = err
			} else {
				c.K = keys
			}

			ch <- c
		}
		close(ch)
	}()
	return ch
}

// ParseConsensus parses a consensus file and extracts the HSDir Hashring
func ParseConsensus(filename string) ([]Hash, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var fingerprint string
	var keys []Hash

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
			var k Hash
			copy(k[len(k)-len(f):], f)
			keys = append(keys, k)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}

type PackedConsensusHdr struct {
	Time Hour
	Len  int32
}

func WritePackedConsensus(w io.Writer, c *Consensus) error {
	hdr := PackedConsensusHdr{Time: c.Time, Len: int32(len(c.K))}
	if err := binary.Write(w, binary.BigEndian, hdr); err != nil {
		return err
	}
	for _, k := range c.K {
		if _, err := w.Write(k[:]); err != nil {
			return err
		}
	}
	return nil
}

type PackReader struct {
	rc  io.ReadCloser
	c   *Consensus
	err error
}

func NewPackReader(filename string) (p *PackReader) {
	f, err := os.Open(filename)
	if err != nil {
		return &PackReader{
			err: err,
		}
	}
	return &PackReader{
		rc: f,
	}
}

func (p *PackReader) Load() bool {
	var hdr PackedConsensusHdr
	if err := binary.Read(p.rc, binary.BigEndian, &hdr); err == io.EOF {
		p.rc.Close()
		return false
	} else if err != nil {
		p.err = err
		p.rc.Close()
		return false
	}
	p.c = &Consensus{
		Time: hdr.Time,
		K:    make([]Hash, hdr.Len),
	}
	for i := int32(0); i < hdr.Len; i++ {
		_, err := io.ReadFull(p.rc, p.c.K[i][:])
		if err != nil {
			p.err = err
			p.rc.Close()
			return false
		}
	}
	return true
}

func (p *PackReader) Consensus() *Consensus {
	return p.c
}

func (p *PackReader) Err() error {
	return p.err
}
