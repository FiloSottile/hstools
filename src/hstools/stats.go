package hstools

import "math/big"

type MetricData struct {
	Mean   *big.Int
	AbsDev *big.Int
	StdDev *big.Int
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

func bigAbsDev(nums []*big.Int, mean *big.Int) *big.Int {
	avg := big.NewInt(0)
	size := big.NewInt(int64(len(nums)))
	for _, n := range nums {
		d := new(big.Int)
		d.Abs(d.Sub(mean, n))
		avg.Add(avg, d.Div(d, size))
	}
	return avg
}

// X versions only consider negative deviations, since attacks only want to be lower

// func bigStdDevX(nums []*big.Int, mean *big.Int) *big.Int {
//     avg := big.NewInt(0)
//     size := big.NewInt(int64(len(nums)) - 1)
//     for _, n := range nums {
//         d := new(big.Int)
//         d.Exp(d.Sub(n, mean), bigTwo, nil)
//         avg.Add(avg, d.Div(d, size))
//     }
//     return bigSqrt(avg)
// }

// func bigAbsDevX(nums []*big.Int, mean *big.Int) *big.Int {
//     avg := big.NewInt(0)
//     size := big.NewInt(int64(len(nums)))
//     for _, n := range nums {
//         d := new(big.Int)
//         d.Abs(d.Sub(mean, n))
//         avg.Add(avg, d.Div(d, size))
//     }
//     return avg
// }

const ROUNDS = 100000

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
		Mean:   mean,
		StdDev: bigStdDev(samples, mean),
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
		Mean:   mean,
		StdDev: bigStdDev(samples, mean),
		AbsDev: bigAbsDev(samples, mean),
	}
}

var mask = new(big.Int).Exp(big.NewInt(2), big.NewInt(160-30), nil)

func Score(v *big.Int, res *MetricData) *big.Int {
	dev := new(big.Int).Sub(res.Mean, v)
	return dev.Div(dev.Mul(dev, big.NewInt(100)), res.AbsDev)
}
