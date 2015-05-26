// +build manually

// lookmeup: bruteforce a onion address that when looked up will have the given
// HSDir, so that you can look it up and see it in the logs of the HSDir
package main

import (
	"bytes"
	"crypto/rand"
	"hstools"
	"log"
	"math/big"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
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
		log.Fatal("usage: lookmeup consensus fingerprint")
	}

	var hash hstools.Hash
	h, err := hstools.FromHex(os.Args[2])
	fatalIfErr(err)
	copy(hash[:], h[:])
	log.Println("[*] Your node is", hstools.ToHex(hash[:]))

	log.Println("[*] Loading consensus...")
	c, err := hstools.ParseConsensus(os.Args[1])
	fatalIfErr(err)
	hashring := hstools.NewHashring(hstools.HashesToIntSlice(c.K))
	hsDirInt := hashring.Prev(new(big.Int).SetBytes(hash[:])) // first HSDir
	hsDirInt = hashring.Prev(hsDirInt)                        // second HSDir
	hsDirInt = hashring.Prev(hsDirInt)                        // third HSDir
	hsDir := hstools.IntToHash(hsDirInt)
	log.Println("    Lowest HSDir:", hstools.ToHex(hsDir[:]))

	log.Println("[*] Starting bruteforce...")
	b := make([]byte, 10)
	for {
		_, err := rand.Read(b)
		fatalIfErr(err)
		onion := hstools.ToBase32(b)
		res, err := hstools.OnionToDescID(onion, time.Now())
		fatalIfErr(err)
		if bytes.Compare(hsDir[:], res[0]) < 0 && bytes.Compare(res[0], hash[:]) < 0 {
			log.Println("    Here, use this:", onion+".onion")
			break
		}
	}

	log.Println("[*] Done!")
}
