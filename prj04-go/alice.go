package main

// Alice the garbler

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
)

type Wire struct {
	// For each wire, we need to prepare two random
	// 128-bit values, v[0] for 0 and v[1] for 1.
	v [2][]byte
}

type Gate struct {
	// NAND, NOR, etc.
	logic string

	// id of input and output wires
	in0, in1, out int

	// garbled truth table
	table [4][]byte
}

// helper function to make gate creation easy
func makeGate(logic string, in0, in1, out int) Gate {
	return Gate{logic: logic, in0: in0, in1: in1, out: out} //assigning values
}

func encryptWires(n, m int) []Wire {
	// We need n+m wires for n input bits and m gates.
	wires := make([]Wire, m+n) //make an array of wire with size m+n

	for i := range wires {
		// use a pointer so we can modify wire in the array
		wire := &wires[i]

		// For each wire, we need to prepare two random
		// 128-bit values, v[0] for 0 and v[1] for 1.
		wire.v[0] = make([]byte, 16)
		wire.v[1] = make([]byte, 16)
		rand.Read(wire.v[0])
		rand.Read(wire.v[1])

		// The first bit is the selection bit. It need
		// to be different for v[0] and v[1] - correct
		// v[1] by inverting its first bit if not.
		if wire.v[0][0]&0x80 == wire.v[1][0]&0x80 {
			wire.v[1][0] = wire.v[1][0] ^ 0x80
		}
	}

	return wires
}

// encrypt a row in the truth table
func encryptOneRow(a, b, o []byte) []byte {
	if len(a) != 16 || len(b) != 16 || len(o) != 16 {
		panic("only 128-bit wires are supported")
	}

	// setup AES block ciphers with 256-bit keys
	c, _ := aes.NewCipher(append(a, b...)) //aes.NewCipher(key []byte) //concatenateing 128 bits of a and 128 bits of b for 256bit key

	// encrypt the output
	garbled := make([]byte, 16)
	c.Encrypt(garbled, o) //encrypting from o to garbled which is empty

	return garbled
}

func encryptGates(gates []Gate, wires []Wire) {
	for i := range gates {
		// use a pointer so we can modify gate in the array
		gate := &gates[i]
		if gate.logic != "NAND" {
			panic("only NAND gates are supported")
		}

		// input and output wires
		A := wires[gate.in0]
		B := wires[gate.in1]
		O := wires[gate.out]

		// selection bits are used to arrange the rows
		// in the garbled truth table
		sa := (A.v[0][0] & 0x80) >> 7 //this gives the first bit either 0 or 1
		sb := (B.v[0][0] & 0x80) >> 7
		println("This is my output")
		println(A.v[0][0])
		println("this is two & 0x80")
		println(A.v[0][0] & 0x80)
		println("this is complete")
		println((A.v[0][0] & 0x80) >> 7)

		// generate rows in the garbled truth table
		gate.table[sa*2+sb] = encryptOneRow(A.v[0], B.v[0], O.v[1]) //if sa=0, sb=1
		gate.table[sa*2+1-sb] = encryptOneRow(A.v[0], B.v[1], O.v[1])
		gate.table[(1-sa)*2+sb] = encryptOneRow(A.v[1], B.v[0], O.v[1])
		gate.table[(1-sa)*2+1-sb] = encryptOneRow(A.v[1], B.v[1], O.v[0])

		fmt.Printf("encrypt %d[%x,%x]=NAND(%d[%x,%x],%d[%x,%x]): %x,%x,%x,%x\n",
			gate.out, O.v[0][:2], O.v[1][:2],
			gate.in0, A.v[0][:2], A.v[1][:2],
			gate.in1, B.v[0][:2], B.v[1][:2],
			gate.table[0][:2], gate.table[1][:2],
			gate.table[2][:2], gate.table[3][:2])
	}
}

func garbleCircuit(gates []Gate) []Wire {
	m := len(gates)
	n := gates[len(gates)-1].out + 1 - m
	wires := encryptWires(n, m)
	encryptGates(gates, wires)
	return wires
}
