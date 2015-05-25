// +build manually

package main

import (
	"hstools"
	"log"
	"math/big"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
)

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if len(os.Args) < 2 {
		log.Fatal("usage: grind pckcns.dat")
	}

	r := hstools.NewPackReader(os.Args[1])
	for r.Load() {
		c := r.Consensus()
		h := hstools.NewHashring(hstools.HashesToIntSlice(c.K))
		log.Println(hstools.HourToTime(c.Time), h.Len())
		res := hstools.SampleDistance4(h)
		log.Printf("%#v", res)

		log.Println("################   MAX", hstools.Score(new(big.Int), res))

		desc, err := hstools.OnionToDescID("silkroadvb5piz3r.onion", hstools.HourToTime(c.Time))
		fatalIfErr(err)
		v := h.Distance4(new(big.Int).SetBytes(desc[0]))
		log.Println("silkroadvb5piz3r.onion", hstools.Score(v, res), v)

		desc, err = hstools.OnionToDescID("facebookcorewwwi.onion", hstools.HourToTime(c.Time))
		fatalIfErr(err)
		v = h.Distance4(new(big.Int).SetBytes(desc[0]))
		log.Println("facebookcorewwwi.onion", hstools.Score(v, res), v)

		log.Println()
	}
	if err := r.Err(); err != nil {
		log.Fatal(err)
	}
}
