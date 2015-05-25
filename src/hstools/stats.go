package hstools

import "math/big"

type MetricData struct {
	Mean   *big.Int
	AbsDev *big.Int
}

type PartitionData struct {
	x0 *big.Int
	x1 *big.Int
	l  *big.Int
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

// func bigAbsDev(nums []*big.Int, mean *big.Int) *big.Int {
//     avg := big.NewInt(0)
//     size := big.NewInt(int64(len(nums)))
//     for _, n := range nums {
//         d := new(big.Int)
//         d.Abs(d.Sub(mean, n))
//         avg.Add(avg, d.Div(d, size))
//     }
//     return avg
// }

func bigAbsDev(nums []*big.Int, mean *big.Int) *big.Int {
	devs := make([]*big.Int, len(nums))
	for i, n := range nums {
		devs[i] = new(big.Int).Sub(n, mean)
		devs[i].Abs(devs[i])
	}
	return bigMean(devs)
}

const ROUNDS = 1000000

func (h *Hashring) Distance(p *big.Int) *big.Int {
	return h.Diff(p, h.Next(p))
}

func SampleDistance(h *Hashring) (res *MetricData) {
	samples := make([]*big.Int, ROUNDS)
	origin := new(big.Int)
	for i := 0; i < ROUNDS; i++ {
		origin.Rand(random, HashringLimit)
		samples[i] = h.Distance(origin)
	}
	mean := bigMean(samples)
	return &MetricData{
		Mean: mean,
		// StdDev: bigStdDev(samples, mean),
		AbsDev: bigAbsDev(samples, mean),
	}
}

func (h *Hashring) Distance4(p *big.Int) *big.Int {
	return h.Diff(p, h.Fourth(p))
}

func SampleDistance4(h *Hashring) (res *MetricData) {
	samples := make([]*big.Int, ROUNDS)
	origin := new(big.Int)
	for i := 0; i < ROUNDS; i++ {
		origin.Rand(random, HashringLimit)
		samples[i] = h.Distance4(origin)
	}
	mean := bigMean(samples)
	return &MetricData{
		Mean: mean,
		// StdDev: bigStdDev(samples, mean),
		AbsDev: bigAbsDev(samples, mean),
	}
}

func (h *Hashring) Distance4Data() (res []*PartitionData) {
	for i, p := range h.points {
		x0 := h.Diff(p, h.points[(i+4)%len(h.points)])
		l := h.Diff(p, h.points[(i+1)%len(h.points)])
		x1 := new(big.Int).Sub(x0, l)
		res = append(res, &PartitionData{
			x0: x0, x1: x1, l: l,
		})
	}
	return
}

func (h *Hashring) DistanceData() (res []*PartitionData) {
	for i, p := range h.points {
		x0 := h.Diff(p, h.points[(i+1)%len(h.points)])
		res = append(res, &PartitionData{
			x0: x0, x1: big.NewInt(0), l: x0,
		})
	}
	return
}

func AnalyzePartitionData(data []*PartitionData) *MetricData {
	m := &MetricData{
		Mean:   new(big.Int),
		AbsDev: new(big.Int),
	}
	for _, part := range data {
		u := new(big.Int).Add(part.x0, part.x1)
		u.Abs(u.Div(u, bigTwo))
		u.Div(u.Mul(u, part.l), HashringLimit)
		m.Mean.Add(m.Mean, u)
	}
	for _, part := range data {
		d0 := new(big.Int).Sub(part.x0, m.Mean)
		d1 := new(big.Int).Sub(part.x1, m.Mean)
		if d0.Sign() == d1.Sign() {
			w := new(big.Int).Add(d0, d1)
			w.Abs(w.Div(w, bigTwo))
			w.Div(w.Mul(w, part.l), HashringLimit)
			m.AbsDev.Add(m.AbsDev, w)
		} else {
			d1.Abs(d1) // Assumes l = x1 - x0 and x1 < x0

			w0 := new(big.Int).Div(d0, bigTwo)
			w0.Div(w0.Mul(w0.Abs(w0), d0), HashringLimit)
			m.AbsDev.Add(m.AbsDev, w0)

			w1 := new(big.Int).Div(d1, bigTwo)
			w1.Div(w1.Mul(w1.Abs(w1), d1), HashringLimit)
			m.AbsDev.Add(m.AbsDev, w1)
		}
	}
	return m
}

var mask = new(big.Int).Exp(big.NewInt(2), big.NewInt(160-30), nil)

func Score(v *big.Int, res *MetricData) *big.Int {
	dev := new(big.Int).Sub(res.Mean, v)
	return dev.Div(dev.Mul(dev, big.NewInt(100)), res.AbsDev)
}
