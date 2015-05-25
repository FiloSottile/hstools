package hstools

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

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
	fing, err := FromHex("EC816FBE76CD94C9064C8F22AF5A468CC46953EA")
	exitIfErr(t, err)
	res := HashIdentity(key.PublicKey)
	if !bytes.Equal(res[:], fing) {
		t.Fail()
	}
}
