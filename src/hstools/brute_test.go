package hstools

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"testing"
)

func exitIfErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestIDDistance(t *testing.T) {
	var a, b, A, B, d IdentityKey
	var aInt, bInt = new(big.Int), new(big.Int)
	var dInt, DInt = new(big.Int), new(big.Int)
	for i := 0; i < 10000; i++ {
		_, err := rand.Read(a[:])
		exitIfErr(t, err)
		copy(A[:], a[:])
		_, err = rand.Read(b[:])
		exitIfErr(t, err)
		copy(B[:], b[:])

		IDDistance(&a, &b, &d)
		DInt = DInt.SetBytes(d[:])

		if !bytes.Equal(a[:], A[:]) || !bytes.Equal(b[:], B[:]) {
			t.Fatal("input changed")
		}

		// Compute the distance by big.Int to check
		aInt = aInt.SetBytes(a[:])
		bInt = bInt.SetBytes(b[:])
		if aInt.Cmp(bInt) < 0 {
			dInt = dInt.Sub(bInt, aInt)
		} else {
			dInt = dInt.Sub(HashringLimit, aInt)
			dInt = dInt.Add(dInt, bInt)
		}

		if DInt.Cmp(dInt) != 0 {
			t.Log(DInt.Bytes())
			t.Log(dInt.Bytes())
			t.Fatal("different result")
		}
	}
}

func TestHashIdentity(t *testing.T) {
	block, _ := pem.Decode([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDX11Z88VBf+4ZJiczyTjTHMS9x1ZbC5qBLQj4LhOWkKJZe9ObK
lcbGd+oyVNip4FTaY5RFenMYOt1ESlYn8jaU/vAi0IMA/E70x9c0p6eLwSr+zCEU
CL/S6ISxwnaYiP92fLfL9keGErKoMbN3t01tAmaDN5jdaaiREVGsHgFVoQIDAQAB
AoGALUw6EHqsfZhR9HkBFBEprmw6Is/KlhjEp0a9srkvYKZL+J25GecZEmn0Mp/v
4Kb9599iLLqoEPu5mC1pq3R/055F97x/IGxxhP/80LmXLCIeeNG+m3s/ezwUNgny
jT+rsCQAxs/r6sjIcCIAfM8rKtXuqcgUew+d8G3hoSwYv+kCQQD9Fc79mdV8sL/c
ChCY9ryxFwofSn8Ljpm4SJ1RssBsXF3+RnG/G6P80k3/wcae/1w1m/KpoqvZT0Qw
fUfMe87XAkEA2lO4P+2oNkjlaqHVlBJYShBm4QoBPls0boX4aB4hjb+AlQ8P024+
Pis7qXa4glxlumlDL6CXQx/cRjsdXyXIRwJBAPEeI/SM6U5Afqm+lQ2GlUMKtkQV
j3CNTXq7A9bgPF+AqLQmnRv704J9Qn6WOQsmMs2IY+ql5p/E2yxvT0ZL9kUCQDkC
bXU8AJWUOVu7wIJ2u9kzKToQG70Foc5Oa0v8ujRCUjgaA77o5ZXkQiMBHjLkH6gq
fmG8ZGMhuaoZG5VRz1cCQGox7SskO48AyaynCKNXM3+vWDNtiwrsxBeX3T2nIWWO
x8IyevfhgPIzX0bajUEqm+phNXWBMUTobyTJbkJQ4NQ=
-----END RSA PRIVATE KEY-----`))
	if block.Type != "RSA PRIVATE KEY" {
		t.Fatal("wrong type")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	exitIfErr(t, err)
	t.Log(key.PublicKey.N.BitLen())
	fing, err := FromHex("EC816FBE76CD94C9064C8F22AF5A468CC46953EA")
	exitIfErr(t, err)
	res := HashIdentity(key.PublicKey)
	if !bytes.Equal(res[:], fing) {
		t.Fail()
	}
}
