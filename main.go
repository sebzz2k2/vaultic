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
		isValidCmd := utils.ValidateCmd(tokens[0])
		if !isValidCmd {
			fmt.Println("Command not implemented")
		}
		isValidArgCount := utils.IsValidArgsCount(tokens[0], len(tokens)-1)
		if !isValidArgCount {
			utils.PrintInvalidArgsError(tokens[0])
		}
		fmt.Println(tokens[0])
	}
}
