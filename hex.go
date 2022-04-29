package go_printer

import (
	"fmt"
	"math"
)

func printHex(data []byte) {
	printHexWithFooter(data, "")
}

func printHexWithFooter(data []byte, footer string) {
	length := len(data)

	for i := 0; i < length; i += 16 {
		fmt.Printf("%0.8x ", i)
		chunkSize := int(math.Min(16, float64(length-i)))
		chunkCap := i + chunkSize

		for j := i; j < chunkCap; j++ {
			if j == i+8 {
				fmt.Print(" ")
			}
			fmt.Printf(" %0.2x", data[j])
		}

		padding := 0
		if chunkSize < 16 {
			padding = (16 - chunkSize) * 3
			if chunkSize < 9 {
				padding++
			}
		}

		for j := 0; j < padding; j++ {
			fmt.Print(" ")
		}

		fmt.Print("  |")

		for j := i; j < chunkCap; j++ {
			if data[j] >= 0x20 && data[j] < 0x7f {
				fmt.Printf("%c", data[j])
			} else {
				fmt.Print(".")
			}
		}

		fmt.Println("|")
	}

	fmt.Printf("%0.8x", len(data))

	if footer != "" {
		fmt.Printf("  %s\n", footer)
	} else {
		fmt.Println()
	}
}
