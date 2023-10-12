package main

import (
	"crypto/aes"
	"fmt"
)

// Bob the evaluator

func evaluateGarbledCircuit(inputs [][]byte, gates []Gate) []byte {
	n := len(inputs) // number of inputs
	m := len(gates)  // number of gates

	// array of signals have a size of n+m
	signals := make([][]byte, m+n)

	// setup inputs signals
	for i := 0; i < n; i++ {
		signals[i] = inputs[i]
		fmt.Printf("input %d=%x\n", i, signals[i][:2])
	}

	// add code below to evaluate the gates

	for _, gate := range gates {
		a := signals[gate.in0]
		b := signals[gate.in1]
		c, _ := aes.NewCipher(append(a, b...)) //using the same key as for encryption in alice.go
		o := make([]byte, 16)                  //making an empty container
		sa := (a[0] & 0x80) >> 7               //this gives the first bit of input a either 0 or 1
		sb := (b[0] & 0x80) >> 7               //this gives the first bit of input b either 0 or 1
		//println("I am printing first bit of a from bob")
		//println(sa)
		//println("I am printing first bit of b from bob")
		//println(sb)
		if sa == 0 && sb == 0 { //using the fisrt bits of the inputs to identify the row and decrypting that exact row.
			c.Decrypt(o, gate.table[0])
		} else if sa == 0 && sb == 1 {
			c.Decrypt(o, gate.table[1])
		} else if sa == 1 && sb == 0 {
			c.Decrypt(o, gate.table[2])
		} else if sa == 1 && sb == 1 {
			c.Decrypt(o, gate.table[3])
		}

		signals[gate.out] = o
	}

	// the last signal is the output
	return signals[n+m-1]
}
