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

	if len(os.Args) < 3 {
		log.Fatal("usage: grind pckcns.dat keys.db")
	}

	keysDB, err := hstools.OpenKeysDb(os.Args[2])
	fatalIfErr(err)

	skipped := 0
	r := hstools.NewPackReader(os.Args[1])
	for r.Load() {
		if skipped < 24*30 {
			skipped++
			continue
		}

		c := r.Consensus()
		h := hstools.NewHashring(hstools.HashesToIntSlice(c.K))
		data, err := h.AgeData(c.Time, keysDB)
		fatalIfErr(err)
		res := hstools.AnalyzePartitionData(data)

		log.Println(hstools.HourToTime(c.Time))
		log.Println("################   MAX", hstools.Score(big.NewInt(0), res))
		log.Println("################   MAX", hstools.Score(big.NewInt(4*24), res))

		desc, err := hstools.OnionToDescID("silkroadvb5piz3r.onion", hstools.HourToTime(c.Time))
		fatalIfErr(err)
		v, err := h.Age(new(big.Int).SetBytes(desc[0]), c.Time, keysDB)
		fatalIfErr(err)
		score := hstools.Score(big.NewInt(int64(v)), res)
		log.Println("AVG", res.Mean, "OUR", v, "SCORE", score)

		desc, err = hstools.OnionToDescID("facebookcorewwwi.onion", hstools.HourToTime(c.Time))
		fatalIfErr(err)
		v, err = h.Age(new(big.Int).SetBytes(desc[0]), c.Time, keysDB)
		fatalIfErr(err)
		score = hstools.Score(big.NewInt(int64(v)), res)
		log.Println("AVG", res.Mean, "OUR", v, "SCORE", score)

		// if score > 230 {
		// 	log.Println(hstools.HourToTime(c.Time), h.Len(), "facebookcorewwwi.onion", score)
		// }

		log.Println()
	}
	if err := r.Err(); err != nil {
		log.Fatal(err)
	}
}
