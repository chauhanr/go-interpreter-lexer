package main

import (
	"os/user"
	"fmt"
	"os"
	"go-interpreter-lexer/repl"
)

func main(){
	user, err := user.Current()
	if err != nil{
		panic (err)
	}
	fmt.Printf("Hello %s! This is Monkey Programming language REPL\n", user.Username)
	fmt.Printf("Type the command for interpretation.")
	repl.Start(os.Stdin, os.Stdout)

}
