package main

import (
	"bytes"
	"fmt"
)

func encryptInputs(inputs []int, wires []Wire) [][]byte {
	signals := make([][]byte, len(inputs))
	for i, input := range inputs {
		signals[i] = wires[i].v[input]
	}
	return signals
}

func decryptOutput(signal []byte, last Wire) int {
	if bytes.Compare(signal, last.v[0]) == 0 {
		return 0
	} else if bytes.Compare(signal, last.v[1]) == 0 {
		return 1
	}
	panic("invalid output signal")
}

func doTest(t string, gates []Gate, inputs []int, expected int) {
	wires := garbleCircuit(gates)
	signals := encryptInputs(inputs, wires)
	signal := evaluateGarbledCircuit(signals, gates)
	output := decryptOutput(signal, wires[len(wires)-1])
	if output != expected {
		panic("incorrect output for test " + t)
	} else {
		fmt.Printf(">>>>>>>>>>>>>>>> test %s pass!\n", t)
	}
}

func test123() {
	gates := []Gate{
		makeGate("NAND", 0, 1, 2),
	}
	doTest("1", gates, []int{0, 0}, 1)
	doTest("2", gates, []int{0, 1}, 1)
	doTest("3", gates, []int{1, 0}, 1)
}

func test4() {
	gates := []Gate{
		makeGate("NAND", 0, 1, 2),
		makeGate("NAND", 2, 2, 3),
	}
	doTest("4", gates, []int{1, 1}, 1)
}

func test567() {
	gates := []Gate{
		makeGate("NAND", 0, 0, 2),
		makeGate("NAND", 1, 1, 3),
		makeGate("NAND", 2, 3, 4),
	}
	doTest("5", gates, []int{0, 0}, 0)
	doTest("6", gates, []int{0, 1}, 1)
	doTest("7", gates, []int{1, 1}, 1)
}

func test890() {
	gates := []Gate{
		makeGate("NAND", 0, 1, 4),
		makeGate("NAND", 2, 3, 5),
		makeGate("NAND", 4, 5, 6),
	}
	doTest("8", gates, []int{0, 0, 0, 0}, 0)
	doTest("9", gates, []int{0, 1, 0, 0}, 0)
	doTest("10", gates, []int{1, 1, 1, 0}, 1)
}

func main() {
	test123()
	test4()
	test567()
	test890()
}
