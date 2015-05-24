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
		return h.points[i].Cmp(p) >= 0
	})
	if i < len(h.points) && h.points[i].Cmp(p) == 0 {
		panic("point is present")
	} else {
		return h.points[i%len(h.points)]
	}
}

func (h *Hashring) Prev(p *big.Int) *big.Int {
	i := sort.Search(len(h.points), func(i int) bool {
		return h.points[i].Cmp(p) >= 0
	})
	return h.points[(i-1)%len(h.points)]
}

func (h *Hashring) Distance(p *big.Int) *big.Int {
	next := h.Next(p)
	if p.Cmp(next) < 0 {
		return new(big.Int).Sub(next, p)
	} else {
		res := new(big.Int).Sub(HashringLimit, p)
		return res.Add(res, next)
	}
}

func bigMean(nums []*big.Int) *big.Int {
	avg := big.NewInt(0)
	size := big.NewInt(int64(len(nums)))
	for _, n := range nums {
		avg.Add(avg, new(big.Int).Div(n, size))
	}
	return avg
}

func bigSqrt(n *big.Int) *big.Int {
	// adapted from mini-gmp
	u, t := new(big.Int), new(big.Int)
	t.SetBit(t, n.BitLen()/2+1, 1)
	for {
		u.Set(t)
		t.Quo(n, u)
		t.Add(t, u)
		t.Rsh(t, 1)
		if t.Cmp(u) >= 0 {
			return u
		}
	}
}

func bigStdDev(nums []*big.Int, mean *big.Int) *big.Int {
	avg := big.NewInt(0)
	size := big.NewInt(int64(len(nums)) - 1)
	for _, n := range nums {
		d := new(big.Int)
		d.Exp(d.Sub(n, mean), bigTwo, nil)
		avg.Add(avg, d.Div(d, size))
	}
	return bigSqrt(avg)
}

func bigMAD(nums []*big.Int, mean *big.Int) *big.Int {
	avg := big.NewInt(0)
	size := big.NewInt(int64(len(nums)))
	for _, n := range nums {
		d := new(big.Int)
		d.Abs(d.Sub(mean, n))
		avg.Add(avg, d.Div(d, size))
	}
	return avg
}

func (h *Hashring) AvgDistance() (mean, stdDev *big.Int) {
	samples := make([]*big.Int, 500000)
	origin := new(big.Int)
	for i := 0; i < 500000; i++ {
		origin.Rand(random, HashringLimit)
		samples[i] = h.Distance(origin)
	}
	mean = bigMean(samples)
	stdDev = bigStdDev(samples, mean)
	return
}
