package main

import (
	"fmt"
	"os"

	"github.com/AlihanE/challenge-json-parser/lexer"
	"github.com/AlihanE/challenge-json-parser/parser"
)

func main() {
	l := lexer.New("{")

	p := parser.New(l)
	_, err := p.ParseProgram()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
