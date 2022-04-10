package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readPassFromInput() string {
	for i := 0; i < 3; i++ {
		if i > 0 {
			fmt.Printf("Try again\n")
		}

		pass1 := readFieldFromInput("SVN password")
		pass2 := readFieldFromInput("SVN password (confirm)")
		if pass1 == pass2 {
			return pass1
		}
		fmt.Printf("Passwords don't match. ")
	}

	fmt.Printf("You so careless. Try later\n")
	os.Exit(1)
	return ""
}

func readFieldFromInput(field string) string {
	fmt.Printf("%s: ", field)
	reader := bufio.NewReader(os.Stdin)

	val, _ := reader.ReadString('\n')
	return strings.Replace(val, "\n", "", -1)
}
