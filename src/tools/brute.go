package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"hstools"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var targetA, targetB hstools.IdentityKey
	if _, err := rand.Read(targetA[:]); err != nil {
		panic(err)
	}
	if _, err := rand.Read(targetB[:]); err != nil {
		panic(err)
	}

	maxA := targetA
	maxA[1] = 0xff
	maxB := targetB
	maxB[1] = 0xff

	log.Println(hstools.ToHex(targetA[:]), hstools.ToHex(maxA[:]))
	log.Println(hstools.ToHex(targetB[:]), hstools.ToHex(maxB[:]))

	keyA, keyB := hstools.Brute(targetA, targetB, maxA, maxB, 3, log.Println)
	if err := pem.Encode(os.Stderr, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyA[0]),
	}); err != nil {
		panic(err)
	}
	if err := pem.Encode(os.Stderr, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(keyB[0]),
	}); err != nil {
		panic(err)
	}
}
