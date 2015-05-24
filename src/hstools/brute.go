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

func Brute(targetA, targetB, maxA, maxB Hash, numKeys, numP int,
	log func(v ...interface{})) (a []*rsa.PrivateKey, b []*rsa.PrivateKey) {
	finished := false
	keys := make(chan *rsa.PrivateKey)
	for p := 0; p < numP; p++ {
		go func(p int) {
			for i := 0; i%100 != 0 || !finished; i++ {
				// We generate real keys because e++ ones are detectable
				// Set the ProbablyPrime rounds to 1 in rand.Prime, we check later
				key, err := rsa.GenerateKey(rand.Reader, 1024)
				if err != nil {
					panic(err)
				}
				id := HashIdentity(key.PublicKey)

				if (bytes.Compare(targetA[:], id[:]) < 0 && bytes.Compare(id[:], maxA[:]) < 0) ||
					(bytes.Compare(targetB[:], id[:]) < 0 && bytes.Compare(id[:], maxB[:]) < 0) {
					keys <- key
				}

				if i%1000 == 0 && i != 0 {
					log("Process #", p, "- iteration #", i)
				}
			}
		}(p)
	}
	for {
		key := <-keys
		if !checkKey(key) {
			log("scrapped bad key")
			continue
		}
		id := HashIdentity(key.PublicKey)
		switch {
		case bytes.Compare(targetA[:], id[:]) < 0 && bytes.Compare(id[:], maxA[:]) < 0:
			a = append(a, key)
		case bytes.Compare(targetB[:], id[:]) < 0 && bytes.Compare(id[:], maxB[:]) < 0:
			b = append(b, key)
		default:
			log("weird, this key is not valid anymore?")
			continue
		}
		log("FOUND ONE!", "A", len(a), "B", len(b))
		if len(a) >= numKeys && len(b) >= numKeys {
			finished = true
			return
		}
	}
}
