// +build manually

package main

import (
	"hstools"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"time"

	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomDistance(ring *hstools.Hashring) *big.Int {
	origin := new(big.Int).Rand(random, hstools.HashringLimit)
	return ring.Distance(origin)
}

func plotDistance() {
	if len(os.Args) < 3 {
		log.Fatal("usage: montecarlo HSDirs RUNS")
	}

	hsdirs, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	runs, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	ring := hstools.RandomHashring(hsdirs)

	v := make(plotter.Values, runs)
	for i := 0; i < runs; i++ {
		d := RandomDistance(ring)
		// keep 32 bits of precision
		v[i] = float64(d.Rsh(d, 160-32).Int64())
	}

	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Title.Text = "Distance"

	h, err := plotter.NewHist(v, 100)
	if err != nil {
		log.Fatal(err)
	}
	p.Add(h)

	if err := p.Save(20, 5, "distance.png"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	plotDistance()
}
