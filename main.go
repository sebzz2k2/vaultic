package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sebzz2k2/vaultic/utils"
)

func main() {
	fmt.Println("Starting vaultic...")
	for {
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		sentence, err := buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		}
		tokens := utils.Tokenize(sentence)

		fmt.Println(tokens[0])
	}
}
