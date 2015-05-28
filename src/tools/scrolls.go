// +build manually

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"hstools"
	"log"
	"math/big"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/mgutz/ansi"
)

var (
	yellow = ansi.ColorFunc("yellow+bh")
	red    = ansi.ColorFunc("red+bh")
	black  = ansi.ColorFunc("red+b:white+h")
)

func fatalIfErr(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}
}

func getAge(h *hstools.Hash, now hstools.Hour, since time.Time, keysDB *hstools.KeysDB) string {
	v, err := keysDB.Lookup(*h)
	fatalIfErr(err)
	if hstools.HourToTime(v.FirstSeen).Before(since) || hstools.HourToTime(v.FirstSeen).Equal(since) {
		return "∞"
	}
	age := now - v.FirstSeen
	switch {
	case age < 24:
		return black(fmt.Sprintf("%d", age))
	case age < 7*24:
		return red(fmt.Sprintf("%d", age))
	case age < 15*24:
		return yellow(fmt.Sprintf("%d", age))
	default:
		return fmt.Sprintf("%d", age)
	}
}

func getLongevity(h *hstools.Hash, now hstools.Hour, until time.Time, keysDB *hstools.KeysDB) string {
	v, err := keysDB.Lookup(*h)
	fatalIfErr(err)
	if hstools.HourToTime(v.LastSeen).After(until) || hstools.HourToTime(v.LastSeen).Equal(until) {
		return "∞"
	}
	long := v.LastSeen - now
	switch {
	case long < 24:
		return black(fmt.Sprintf("%d", long))
	case long < 7*24:
		return red(fmt.Sprintf("%d", long))
	case long < 15*24:
		return yellow(fmt.Sprintf("%d", long))
	default:
		return fmt.Sprintf("%d", long)
	}
}

func getColo(h *hstools.Hash, keysDB *hstools.KeysDB) string {
	res, _ := hstools.ColocatedKeys(h[:], keysDB)
	switch {
	case res > 20:
		return black(fmt.Sprintf("%d", res))
	case res > 10:
		return red(fmt.Sprintf("%d", res))
	case res > 5:
		return yellow(fmt.Sprintf("%d", res))
	default:
		return fmt.Sprintf("%d", res)
	}
}

func colorScore(v int64) string {
	switch {
	case v > 200:
		return black(fmt.Sprintf("%d", v))
	case v > 150:
		return red(fmt.Sprintf("%d", v))
	case v > 100:
		return yellow(fmt.Sprintf("%d", v))
	default:
		return fmt.Sprintf("%d", v)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.SetFlags(0)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if len(os.Args) != 7 {
		log.Fatal("usage: scavenge pckcns.dat keys.db stats.jsonl xxx.onion 2015-05-01-01 2015-05-31-23")
	}

	keysDB, err := hstools.OpenKeysDb(os.Args[2])
	fatalIfErr(err)

	since, err := time.Parse("2006-01-02-15", os.Args[5])
	if err != nil {
		log.Fatal(err)
	}
	until, err := time.Parse("2006-01-02-15", os.Args[6])
	if err != nil {
		log.Fatal(err)
	}

	statsFile, err := os.Open(os.Args[3])
	fatalIfErr(err)
	jsonStats := json.NewDecoder(bufio.NewReaderSize(statsFile, 1024*1024))
	var metrics hstools.AnalyzedConsensus
	var last0, last1 []byte

	r := hstools.NewPackReader(os.Args[1])
	for r.Load() {
		c := r.Consensus()
		jsonStats.Decode(&metrics)
		if metrics.T != c.Time {
			log.Fatal("unaligned data")
		}

		if since.After(hstools.HourToTime(c.Time)) {
			continue
		}
		if until.Before(hstools.HourToTime(c.Time)) {
			break
		}

		h := hstools.NewHashring(hstools.HashesToIntSlice(c.K))

		desc, err := hstools.OnionToDescID(os.Args[4], hstools.HourToTime(c.Time))
		fatalIfErr(err)
		desc0 := new(big.Int).SetBytes(desc[0])
		desc1 := new(big.Int).SetBytes(desc[1])

		hsDir0 := hstools.IntsToHashSlice(h.Next3(desc0))
		allHSDir0 := append(append(append([]byte{}, hsDir0[0][:]...), hsDir0[1][:]...), hsDir0[2][:]...)
		hsDir1 := hstools.IntsToHashSlice(h.Next3(desc1))
		allHSDir1 := append(append(append([]byte{}, hsDir1[0][:]...), hsDir1[1][:]...), hsDir1[2][:]...)
		if bytes.Equal(last0, allHSDir0) && bytes.Equal(last1, allHSDir1) {
			continue
		}
		last0, last1 = allHSDir0, allHSDir1

		// Metric 1. Age
		// age0, err := h.Age(desc0, c.Time, keysDB)
		// fatalIfErr(err)
		// age1, err := h.Age(desc1, c.Time, keysDB)
		// fatalIfErr(err)

		// Metric 2. Longevity
		// long0, err := h.Longevity(desc0, c.Time, keysDB)
		// fatalIfErr(err)
		// long1, err := h.Longevity(desc1, c.Time, keysDB)
		// fatalIfErr(err)

		// Metric 3. Distance
		dist0 := hstools.Score(h.Distance(desc0), metrics.Distance)
		dist1 := hstools.Score(h.Distance(desc1), metrics.Distance)

		// Metric 4. Distance4
		dist40 := hstools.Score(h.Distance4(desc0), metrics.Distance4)
		dist41 := hstools.Score(h.Distance4(desc1), metrics.Distance4)

		// Metric 5. Colocated keys
		// colo0 := h.Colocated(desc0, keysDB)
		// colo1 := h.Colocated(desc1, keysDB)

		fmt.Printf(`
###### %s
###### Replica 0 - Dist score %s - Dist4 score %s
%s - Age %s - Long %s - Colo keys %s
%s - Age %s - Long %s - Colo keys %s
%s - Age %s - Long %s - Colo keys %s
###### Replica 1 - Dist score %s - Dist4 score %s
%s - Age %s - Long %s - Colo keys %s
%s - Age %s - Long %s - Colo keys %s
%s - Age %s - Long %s - Colo keys %s
`,
			hstools.HourToTime(c.Time), colorScore(dist0), colorScore(dist40),

			hstools.ToHex(hsDir0[0][:]), getAge(hsDir0[0], c.Time, since, keysDB),
			getLongevity(hsDir0[0], c.Time, until, keysDB), getColo(hsDir0[0], keysDB),

			hstools.ToHex(hsDir0[1][:]), getAge(hsDir0[1], c.Time, since, keysDB),
			getLongevity(hsDir0[1], c.Time, until, keysDB), getColo(hsDir0[1], keysDB),

			hstools.ToHex(hsDir0[2][:]), getAge(hsDir0[2], c.Time, since, keysDB),
			getLongevity(hsDir0[2], c.Time, until, keysDB), getColo(hsDir0[2], keysDB),

			colorScore(dist1), colorScore(dist41),

			hstools.ToHex(hsDir1[0][:]), getAge(hsDir1[0], c.Time, since, keysDB),
			getLongevity(hsDir1[0], c.Time, until, keysDB), getColo(hsDir1[0], keysDB),

			hstools.ToHex(hsDir1[1][:]), getAge(hsDir1[1], c.Time, since, keysDB),
			getLongevity(hsDir1[1], c.Time, until, keysDB), getColo(hsDir1[1], keysDB),

			hstools.ToHex(hsDir1[2][:]), getAge(hsDir1[2], c.Time, since, keysDB),
			getLongevity(hsDir1[2], c.Time, until, keysDB), getColo(hsDir1[2], keysDB),
		)

		// if colo0 > 20 || colo1 > 20 {
		// 	log.Println(hstools.HourToTime(c.Time), colo0, colo1)
		// }

		// if dist0 > 200 || dist1 > 200 || dist40 > 250 || dist41 > 250 ||
		// 	age0 < 24*7 || age1 < 24*7 || long0 < 24*3 || long1 < 24*3 ||
		// 	colo0 > 20 || colo1 > 20 {
		// 	log.Println(hstools.HourToTime(c.Time),
		// 		"# AGE", age0, age1,
		// 		"# LONG", long0, long1,
		// 		"# DIST", dist0, dist1,
		// 		"# DIST4", dist40, dist41,
		// 		"# COLO", colo0, colo1)
		// }
	}
	fatalIfErr(r.Err())
}
