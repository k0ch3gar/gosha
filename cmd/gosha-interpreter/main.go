package main

import (
	"fmt"
	"kstmc.com/gosha/internal/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Дарова %s!\nЭто Гоша. Он очень крутой!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
