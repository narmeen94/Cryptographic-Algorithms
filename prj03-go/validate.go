// crypto/crypto.go
package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func validateRSAOAEP() {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := &privateKey.PublicKey

	plaintxt := "This is a secret!"
	label := ""

	ciphertxt, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader,
		publicKey, []byte(plaintxt), []byte(label))

	pbuf, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader,
		privateKey, ciphertxt, []byte(label))
	plaintxt2 := string(pbuf)

	fmt.Printf("Validated %t: RSA-OAEP(%s)=%x\n",
		plaintxt == plaintxt2, plaintxt, ciphertxt)
}

func validateRSAPSS() {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	publicKey := &privateKey.PublicKey

	msg := "Hello world!"
	hash := sha256.Sum256([]byte(msg))

	signature, _ := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hash[:], nil)

	err := rsa.VerifyPSS(publicKey, crypto.SHA256, hash[:], signature, nil)

	fmt.Printf("Validated %t: RSA-PSS(SHA256(%s)=%x)=%x\n",
		err == nil, msg, hash, signature)
}
