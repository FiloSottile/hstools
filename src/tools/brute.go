// +build manually

// brute: bruteforce Identity Keys that will be the 6 HSDir for the given onion
// at the given time, considering the given consensus state
package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

	if len(os.Args) < 4 {
		log.Fatal("usage: brute consensus onion RFC3339time")
	}

	log.Println("[*] Computing HS descriptors...")
	t, err := time.Parse(time.RFC3339, os.Args[3])
	fatalIfErr(err)
	desc, err := hstools.OnionToDescID(os.Args[2], t)
	fatalIfErr(err)
	var keyA, keyB hstools.Hash
	copy(keyA[:], desc[0])
	copy(keyB[:], desc[1])
	log.Printf("    Onion '%s' at time '%s'\n", os.Args[2], t)
	log.Println("    Descriptor A:", hstools.ToHex(keyA[:]))
	log.Println("    Descriptor B:", hstools.ToHex(keyB[:]))

	log.Println("[*] Loading consensus...")
	// Note: this should maybe also consider potential future HSDir
	c, err := hstools.ParseConsensus(os.Args[1])
	fatalIfErr(err)
	hashring := hstools.NewHashring(hstools.HashesToIntSlice(c.K))
	nextA := hstools.IntToHash(hashring.Next(new(big.Int).SetBytes(keyA[:])))
	nextB := hstools.IntToHash(hashring.Next(new(big.Int).SetBytes(keyB[:])))
	log.Println("    First HSDir A:", hstools.ToHex(nextA[:]))
	log.Println("    First HSDir B:", hstools.ToHex(nextB[:]))

	log.Println("[*] Starting bruteforce...")
	keysA, keysB := hstools.Brute(keyA, keyB, nextA, nextB, 3,
		runtime.NumCPU(), log.Println)

	log.Println("[*] Done!")

	for _, keys := range [][]*rsa.PrivateKey{keysA, keysB} {
		for i := 0; i < 3; i++ {
			if err := pem.Encode(os.Stderr, &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(keys[i]),
			}); err != nil {
				log.Fatal(err)
			}
		}
	}
}
