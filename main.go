package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlihanE/challenge-json-parser/ast"
	"github.com/AlihanE/challenge-json-parser/lexer"
	"github.com/AlihanE/challenge-json-parser/parser"
)

func main() {
	l := lexer.New("{")

	p := parser.New(l)
	v, err := p.ParseProgram()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Println("Result", v.Type)
	r := (*v.RootValue).(ast.Object)
	log.Println(r)
}
