// crypto/crypto.go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func findPassword() {
	nonce := make([]byte, 12)
	data := "jwang34@iit.edu"

	ciphermac, _ := hex.DecodeString(
		"7d4eb640daf43844bd605bb1c2a66046fd33cba7d2ba35828d25056c953834d93b9a04c54fa86147f62a")

	fmt.Printf("finding password for nonce=%x, data=%s, ciphermac=%x...\n",
		nonce, data, ciphermac) //i.e key is missing. if key is found, we can decrypt by using AES

	//Brute-Force Attack. i.e. trying all combinations for 4 digit number so 10000 combinations.
	for i := 0; i < 10000; i++ {
		password := fmt.Sprintf("%04d", i)

		// Modify code below to check that if sha256(password) is a correct
		// key to decrypt the ciphertext with MAC in ciphermac, which is
		// obtained via AES-GCM with a 0 nouce and the data being jwang34@iit.edu
		//
		// Once found, show the correct password as well as the plaintext.

		//creating a hash function of the password using sha-256
		sha := sha256.New()
		sha.Write([]byte(password)) //this converts the string to bytes.

		hash := sha.Sum(nil) //this is the final hash
		key := hash          //this hash behaves as a key for the AES-GCM.

		//initiating a block for AES and giving it a key
		block, _ := aes.NewCipher(key)
		aesgcm, _ := cipher.NewGCM(block)

		//decryption by AES (AES takes in 4 parameters for decryption. a key,nonce,chiphertext and AAD)
		//pbuf is buffer for the plaintext i.e. result of decryption. err is the error.
		pbuf, err := aesgcm.Open(nil, nonce, ciphermac, []byte(data))

		//as soon as the correct combination of 4 digits is found, the loop is break and the printing is done.
		if err == nil {
			fmt.Printf("Validated %t: AES-GCM(%s, data=%s, nonce=%x, key=%x,password=%s)=%x\n",
				err == nil, string(pbuf), data, nonce, key, password, ciphermac)
			fmt.Printf("  correct password=%s\n", password)
			break

		}

	}
}

func main() {
	//validateSHA256()
	//validateAESGCM()
	findPassword()
}
