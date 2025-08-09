package main

import (
	"fmt"
	"os"
	"os/user"

	"kstmc.com/gosha/internal/repl"
	"kstmc.com/gosha/internal/token"
)

func main() {
	err := token.SetupBashCalls()
	if err != nil {
		return
	}

	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err)
			os.Exit(1)
		}

		defer file.Close()
		repl.Start(file, os.Stdout)
	} else {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Hi %s!\nThat's Gosha!\n", user.Username)
		repl.Start(os.Stdin, os.Stdout)
	}
}
