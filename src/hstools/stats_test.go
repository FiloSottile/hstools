package hstools

import "testing"

func TestAnalyzePartitionData(t *testing.T) {
	h := RandomHashring(3000)
	res := AnalyzePartitionData(h.Distance4Data())
	t.Log(res)

	var meanHi, meanLo, devHi, devLo bool
	for i := 0; i < 10; i++ {
		resSample := SampleDistance4(h)
		t.Log(resSample)
		meanHi = meanHi || resSample.Mean.Cmp(res.Mean) > 0
		meanLo = meanLo || resSample.Mean.Cmp(res.Mean) < 0
		devHi = devHi || resSample.AbsDev.Cmp(res.AbsDev) > 0
		devLo = devLo || resSample.AbsDev.Cmp(res.AbsDev) < 0
	}

	if !(meanHi && meanLo && devHi && devLo) {
		t.Failed()
	}
}
