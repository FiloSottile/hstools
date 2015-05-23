package hstools

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/asn1"
	"math"
)

// IDDistance efficiently calculates the difference (b - a) mod 2^160, or
// distance a -> b on a 20-byte ring and stores it in d. a and b unchanged.
func IDDistance(a, b, d *Hash) {
	var carry bool
	for i := len(a) - 1; i >= 0; i-- {
		B := b[i]
		if carry {
			B--
		}
		d[i] = B - a[i]
		carry = B < a[i] || (carry && B == math.MaxUint8)
	}
}

func HashIdentity(pk rsa.PublicKey) Hash {
	// tor-spec.txt#n108
	// When we refer to "the hash of a public key", we mean the SHA-1 hash of the
	// DER encoding of an ASN.1 RSA public key (as specified in PKCS.1).
	// rfc3447#appendix-A.1.1
	// RSAPublicKey ::= SEQUENCE {
	//     modulus           INTEGER,  -- n
	//     publicExponent    INTEGER   -- e
	// }
	der, err := asn1.Marshal(pk)
	if err != nil {
		panic(err)
	}
	return Hash(sha1.Sum(der))
}

func checkKey(key *rsa.PrivateKey) bool {
	for _, p := range key.Primes {
		if !p.ProbablyPrime(20) {
			return false
		}
	}
	return true
}

func Brute(targetA, targetB, maxA, maxB Hash, n int,
	log func(v ...interface{})) (a []*rsa.PrivateKey, b []*rsa.PrivateKey) {
	for i := 0; ; i++ {
		// We generate real keys because e++ ones are detectable
		// Set the ProbablyPrime rounds to 1 in rand.Prime, we check later
		key, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			panic(err)
		}
		id := HashIdentity(key.PublicKey)

		if bytes.Compare(targetA[:], id[:]) < 0 &&
			bytes.Compare(id[:], maxA[:]) < 0 {
			if !checkKey(key) {
				log("scrapped bad key", i)
				continue
			}
			a = append(a, key)
			log(len(a), len(b), i)
			if len(a) >= n && len(b) >= n {
				return
			}
		}

		if bytes.Compare(targetB[:], id[:]) < 0 &&
			bytes.Compare(id[:], maxB[:]) < 0 {
			if !checkKey(key) {
				log("scrapped bad key", i)
				continue
			}
			b = append(b, key)
			log(len(a), len(b), i)
			if len(a) >= n && len(b) >= n {
				return
			}
		}

		if i%1000 == 0 && i != 0 {
			log(len(a), len(b), i)
		}
	}
}
