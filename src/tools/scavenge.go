// +build manually

package main

import (
	"bufio"
	"encoding/json"
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

	if len(os.Args) != 5 {
		log.Fatal("usage: scavenge pckcns.dat keys.db stats.jsonl xxx.onion")
	}

	keysDB, err := hstools.OpenKeysDb(os.Args[2])
	fatalIfErr(err)

	statsFile, err := os.Open(os.Args[3])
	fatalIfErr(err)
	jsonStats := json.NewDecoder(bufio.NewReaderSize(statsFile, 1024*1024))
	var metrics hstools.AnalyzedConsensus

	skipped := 0
	r := hstools.NewPackReader(os.Args[1])
	for r.Load() {
		c := r.Consensus()
		jsonStats.Decode(&metrics)
		if metrics.T != c.Time {
			log.Fatal("unaligned data")
		}

		if skipped < 24*30 {
			skipped++
			continue
		}

		h := hstools.NewHashring(hstools.HashesToIntSlice(c.K))

		desc, err := hstools.OnionToDescID(os.Args[4], hstools.HourToTime(c.Time))
		fatalIfErr(err)
		desc0 := new(big.Int).SetBytes(desc[0])
		desc1 := new(big.Int).SetBytes(desc[1])

		// Metric 1. Age
		age0, err := h.Age(desc0, c.Time, keysDB)
		fatalIfErr(err)
		age1, err := h.Age(desc1, c.Time, keysDB)
		fatalIfErr(err)

		// Metric 2. Longevity
		long0, err := h.Longevity(desc0, c.Time, keysDB)
		fatalIfErr(err)
		long1, err := h.Longevity(desc1, c.Time, keysDB)
		fatalIfErr(err)

		// Metric 3. Distance
		dist0 := hstools.Score(h.Distance(desc0), metrics.Distance)
		dist1 := hstools.Score(h.Distance(desc1), metrics.Distance)

		// Metric 4. Distance4
		dist40 := hstools.Score(h.Distance4(desc0), metrics.Distance4)
		dist41 := hstools.Score(h.Distance4(desc1), metrics.Distance4)

		// Metric 5. Colocated keys
		colo0 := h.Colocated(desc0, keysDB)
		colo1 := h.Colocated(desc1, keysDB)

		if colo0 > 20 || colo1 > 20 {
			log.Println(hstools.HourToTime(c.Time), colo0, colo1)
		}

		if dist0 > 200 || dist1 > 200 || dist40 > 250 || dist41 > 250 ||
			age0 < 24*7 || age1 < 24*7 || long0 < 24*3 || long1 < 24*3 ||
			colo0 > 20 || colo1 > 20 {
			log.Println(hstools.HourToTime(c.Time),
				"# AGE", age0, age1,
				"# LONG", long0, long1,
				"# DIST", dist0, dist1,
				"# DIST4", dist40, dist41,
				"# COLO", colo0, colo1)
		}
	}
	if err := r.Err(); err != nil {
		log.Fatal(err)
	}
}
