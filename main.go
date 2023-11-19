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
	l := lexer.New(`[{}]`)

	p := parser.New(l)
	v, err := p.ParseProgram()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(p.Errors()) > 0 {
		panic(p.Errors())
	}

	log.Println("Result", v.Type)
	switch (*v.RootValue).(type) {
	case ast.Object:
		r := (*v.RootValue).(ast.Object)
		log.Println(r)
	case ast.Array:
		r := (*v.RootValue).(ast.Array)
		log.Println(r)
	}
}
