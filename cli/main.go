package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/sebzz2k2/vaultic/internal/protocol/lexer"
	"github.com/sebzz2k2/vaultic/internal/resp"
)

func main() {
	var host = flag.String("host", "localhost", "Vaultic server host")
	var port = flag.String("port", "5381", "Vaultic server port")
	flag.Parse()

	// Establish TCP connection
	address := fmt.Sprintf("%s:%s", *host, *port)
	fmt.Printf("Connecting to Vaultic at %s\n", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Error: Failed to connect to Vaultic server at %s: %v\n", address, err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to Vaultic server!")
	fmt.Println("Hello, Vaultic!")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("vaultic> ")
		if !scanner.Scan() {
			break
		}
		tknStr := lexer.TokenizeCLI(scanner.Text())

		// Send the tokenized string to the server
		_, err := conn.Write([]byte(tknStr))
		if err != nil {
			fmt.Printf("Error sending data to server: %v\n", err)
			continue
		}

		decodedResponse, err := resp.NewDecoder(bufio.NewReader(conn)).Decode()
		if err != nil {
			fmt.Printf("Error reading response from server: %v\n", err)
			continue
		}
		for _, v := range decodedResponse.Array {
			fmt.Println(v.String)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}
}
