// crypto/crypto.go
package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func unused(i interface{}) {}

func genEncKeyAndSign(auth *rsa.PrivateKey) (
	enc *rsa.PrivateKey, encPubSig []byte) {
	enc, _ = rsa.GenerateKey(rand.Reader, 1024)

	msg := fmt.Sprintf("N=%v E=%v", enc.PublicKey.N, enc.PublicKey.E)
	hash := sha256.Sum256([]byte(msg))

	encPubSig, _ = rsa.SignPSS(rand.Reader, auth, crypto.SHA256, hash[:], nil)

	return
}

func aliceVerify(bobAuthPub, bobEncPub *rsa.PublicKey, bobEncPubSig []byte) bool {
	msg := fmt.Sprintf("N=%v E=%v", bobEncPub.N, bobEncPub.E)
	hash := sha256.Sum256([]byte(msg))

	// Modify code below to check if the signature of hash is bobEncPubSig
	// using bobAuthPub. Return true if yes, or false otherwise.
	err := rsa.VerifyPSS(bobAuthPub, crypto.SHA256, hash[:], bobEncPubSig, nil)

	if err == nil {
		fmt.Printf("verified by alice: %t\n", err == nil)
		return true

	}

	return false

	//unused(hash)

}

func aliceEncrypt(bobEncPub *rsa.PublicKey, plaintxt []byte) []byte {

	// Modify code below to encrypt plaintxt using bobEncPub
	// and return the ciphertxt.
	label := ""
	ciphertxt, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, bobEncPub, []byte(plaintxt), []byte(label))
	return ciphertxt
}

func bobDecrypt(bobEnc *rsa.PrivateKey, ciphertxt []byte) []byte {
	// Modify code below to decrypt ciphertxt using bobEnc
	// and return the plaintxt.
	label := ""

	pbuf, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader, bobEnc, ciphertxt, []byte(label))
	plaintxt := string(pbuf)

	return []byte(plaintxt)
}

func doubleRSAPFS() {
	bobAuth, _ := rsa.GenerateKey(rand.Reader, 1024)

	plaintxt := "Hello Bob this is Alice."

	bobEnc, bobEncSig := genEncKeyAndSign(bobAuth)

	if !aliceVerify(&bobAuth.PublicKey, &bobEnc.PublicKey, bobEncSig) {
		panic("Man-in-the-middle attack detected")
	}

	ciphertxt := aliceEncrypt(&bobEnc.PublicKey, []byte(plaintxt))
	pbuf := bobDecrypt(bobEnc, ciphertxt)
	if plaintxt != string(pbuf) {
		panic("decrypt fails")
	}
	fmt.Println("doubleRSAPFS completed successfully.")
}

func doubleRSAPFSMITM() {
	bobAuth, _ := rsa.GenerateKey(rand.Reader, 1024)
	oscarAuth, _ := rsa.GenerateKey(rand.Reader, 1024)

	oscarEnc, oscarEncSig := genEncKeyAndSign(oscarAuth)

	if aliceVerify(&bobAuth.PublicKey, &oscarEnc.PublicKey, oscarEncSig) {
		panic("Man-in-the-middle attack not detected")
	}

	fmt.Println("doubleRSAPFSMITM completed successfully.")
}

func main() {
	//validateRSAOAEP()
	//validateRSAPSS()
	doubleRSAPFS()
	doubleRSAPFSMITM()
}
