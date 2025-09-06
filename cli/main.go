package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sebzz2k2/vaultic/internal/protocol/lexer"
	"github.com/sebzz2k2/vaultic/internal/resp"
)

func main() {
	fmt.Println("Hello, Vaultic!")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		tknStr := lexer.TokenizeCLI(scanner.Text())
		result, err := resp.DecodeString(tknStr)
		if err != nil {
			fmt.Println("Error decoding response:", err)
			continue
		}

		if result.Type != "array" {
			fmt.Println("Error: Expected array response")
			continue
		}

		if len(result.Array) == 0 {
			fmt.Println("Error: Empty command")
			continue
		}

		fmt.Printf("Command: %s\n", result.Array[0].String)

		for i := 1; i < len(result.Array); i++ {
			fmt.Printf("Arg[%d]: %s\n", i, result.Array[i].String)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
