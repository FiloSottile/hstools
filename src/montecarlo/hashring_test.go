package main

import (
	"math/big"
	"testing"
)

func TestHashringNext(t *testing.T) {
	h := NewHashring([]*big.Int{
		big.NewInt(500),
		big.NewInt(200),
		big.NewInt(300),
		big.NewInt(400),
		big.NewInt(100),
	})

	for i := int64(1); i <= 5; i++ {
		expected := i * 100
		n := h.Next(big.NewInt(expected - 50))
		if n.Int64() != expected {
			t.Fatal(n.Int64(), expected)
		}
	}

	n := h.Next(big.NewInt(550))
	if n.Int64() != 100 {
		t.Fatal(n.Int64(), 100)
	}
}

func TestHashringDistance(t *testing.T) {
	h := NewHashring([]*big.Int{
		big.NewInt(500),
		big.NewInt(200),
		big.NewInt(300),
		big.NewInt(400),
		big.NewInt(100),
	})

	n := h.Distance(big.NewInt(450))
	if n.Int64() != 50 {
		t.Fatal(n.Int64())
	}

	n = h.Distance(big.NewInt(550))
	exp := new(big.Int).Sub(HashringLimit, big.NewInt(550-100))
	if n.Cmp(exp) != 0 {
		t.Fatal(n, exp)
	}
}

func TestAvgDistance(t *testing.T) {
	ring := RandomHashring(2500)
	mean, stdDev := ring.AvgDistance()
	t.Log(mean, stdDev)
}
