// +build manually

// grind: compute the average distance to 1st and 4th node and its average
// deviation and output it as JSON lines
package main

import (
	"bufio"
	"encoding/json"
	"hstools"
	"log"
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

	if len(os.Args) != 2 {
		log.Fatal("usage: grind pckcns.dat > stats.jsonl")
	}

	w := bufio.NewWriterSize(os.Stdout, 1024*1024)
	jsonEncoder := json.NewEncoder(w)
	r := hstools.NewPackReader(os.Args[1])
	for i := 0; r.Load(); i++ {
		c := r.Consensus()
		h := hstools.NewHashring(hstools.HashesToIntSlice(c.K))

		fatalIfErr(jsonEncoder.Encode(hstools.AnalyzedConsensus{
			T:         c.Time,
			Distance:  hstools.AnalyzePartitionData(h.DistanceData()),
			Distance4: hstools.AnalyzePartitionData(h.Distance4Data()),
		}))

		if i%1000 == 0 {
			log.Println(hstools.HourToTime(c.Time))
		}
	}
	fatalIfErr(r.Err())
	fatalIfErr(w.Flush())
}
