package cliutil

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Confirm(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s (yes/no): ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			continue
		}
	}
}

func Run(v string, fn func() error) error {
	fmt.Print("[ ] ", v)
	if err := fn(); err != nil {
		fmt.Println("\r[-]", v)
		return err
	}

	fmt.Println("\r[+]", v)

	return nil
}

func Error(v string) {
	fmt.Printf("[-] %s\n", v)
}
