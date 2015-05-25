package hstools

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
)

func exitIfErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestHashDistance(t *testing.T) {
	var a, b, A, B, d Hash
	var aInt, bInt = new(big.Int), new(big.Int)
	var dInt, DInt = new(big.Int), new(big.Int)
	for i := 0; i < 10000; i++ {
		_, err := rand.Read(a[:])
		exitIfErr(t, err)
		copy(A[:], a[:])
		_, err = rand.Read(b[:])
		exitIfErr(t, err)
		copy(B[:], b[:])

		a.Distance(&b, &d)
		DInt = DInt.SetBytes(d[:])

		if !bytes.Equal(a[:], A[:]) || !bytes.Equal(b[:], B[:]) {
			t.Fatal("input changed")
		}

		// Compute the distance by big.Int to check
		aInt = aInt.SetBytes(a[:])
		bInt = bInt.SetBytes(b[:])
		if aInt.Cmp(bInt) < 0 {
			dInt = dInt.Sub(bInt, aInt)
		} else {
			dInt = dInt.Sub(HashringLimit, aInt)
			dInt = dInt.Add(dInt, bInt)
		}

		if DInt.Cmp(dInt) != 0 {
			t.Log(DInt.Bytes())
			t.Log(dInt.Bytes())
			t.Fatal("different result")
		}
	}
}
