package hstools

import (
	"math/big"
	"math/rand"
	"sort"
	"time"
)

var (
	bigOne = big.NewInt(1)
	bigTwo = big.NewInt(2)

	HashringLimit = new(big.Int).Lsh(bigOne, 160)

	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Hashring struct {
	// points is the sorted list of values present on the ring
	points []*big.Int
}

type bigIntSlice []*big.Int

func (b bigIntSlice) Len() int           { return len(b) }
func (b bigIntSlice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b bigIntSlice) Less(i, j int) bool { return b[i].Cmp(b[j]) < 0 }

// NewHashring returns a Hashring with the given unsorted points.
func NewHashring(points []*big.Int) *Hashring {
	h := &Hashring{
		points: make([]*big.Int, len(points)),
	}
	copy(h.points, points)
	sort.Sort(bigIntSlice(h.points))
	return h
}

func RandomHashring(entries int) *Hashring {
	points := make([]*big.Int, entries)
	for i := 0; i < entries; i++ {
		points[i] = new(big.Int).Rand(random, HashringLimit)
	}

	return NewHashring(points)
}

func (h *Hashring) Len() int {
	return len(h.points)
}

func (h *Hashring) Next(p *big.Int) *big.Int {
	i := sort.Search(len(h.points), func(i int) bool {
		return h.points[i].Cmp(p) > 0
	})
	return h.points[i%len(h.points)]
}

func (h *Hashring) Fourth(p *big.Int) *big.Int {
	i := sort.Search(len(h.points), func(i int) bool {
		return h.points[i].Cmp(p) > 0
	})
	return h.points[(i+3)%len(h.points)]
}

func (h *Hashring) Prev(p *big.Int) *big.Int {
	i := sort.Search(len(h.points), func(i int) bool {
		return h.points[i].Cmp(p) >= 0
	})
	return h.points[(i-1)%len(h.points)]
}

func (*Hashring) Diff(from, to *big.Int) *big.Int {
	if from.Cmp(to) < 0 {
		return new(big.Int).Sub(to, from)
	} else {
		res := new(big.Int).Sub(HashringLimit, from)
		return res.Add(res, to)
	}
}
