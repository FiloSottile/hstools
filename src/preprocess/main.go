package main

import (
	"hstools"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("usage: preprocess /data/dir/ 2014-01-01-00 2014-01-31-23")
	}

	pckFile, err := os.Create("pckcns.dat")
	if err != nil {
		log.Fatal(err)
	}
	keysDB, err := hstools.OpenKeysDb("keys.db")
	if err != nil {
		log.Fatal(err)
	}

	since, err := time.Parse("2006-01-02-15", os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	until, err := time.Parse("2006-01-02-15", os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	pprofF, err := os.Create("preprocess.pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(pprofF)

	ch := hstools.ReadConsensuses(os.Args[1],
		hstools.Hour(since.Unix()/3600), hstools.Hour(until.Unix()/3600))
	for c := range ch {
		if c.Error != nil {
			log.Println(c.Error)
			continue
		}

		// desc, err := hstools.OnionToDescID("facebookcorewwwi.onion",
		// 	time.Unix(int64(c.Time*3600), 0))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// d := hstools.ToHex(c.H.Distance(new(big.Int).SetBytes(desc[0])).Bytes())
		// n := hstools.ToHex(c.H.Next(new(big.Int).SetBytes(desc[0])).Bytes())
		// strings.Repeat("0", 40-len(d))+d

		if err := hstools.WritePackedConsensus(pckFile, c); err != nil {
			log.Fatal(err)
		}

		if err := keysDB.Seen(c.K, c.Time); err != nil {
			log.Fatal(err)
		}

		log.Println(c.Filename, len(c.K))
	}

	pprof.StopCPUProfile()

	if err := keysDB.Close(); err != nil {
		log.Fatal(err)
	}
}