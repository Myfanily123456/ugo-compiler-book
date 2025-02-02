//go:build ignore

package main

import (
	"fmt"
	"os"

	lexpkg "github.com/chai2010/ugo/lexer"
)

func main() {
	code := loadCode("../hello.ugo")
	lexer := lexpkg.NewLexer("../hello.ugo", code)

	for i, tok := range lexer.Tokens() {
		fmt.Printf(
			"%02d: %-12v: %-20q // %s\n",
			i, tok.Type, tok.Literal,
			tok.Pos.Position("../hello.ugo", code),
		)
	}

	fmt.Println("----")

	for i, tok := range lexer.Comments() {
		fmt.Printf(
			"%02d: %-12v: %-20q // %s\n",
			i, tok.Type, tok.Literal,
			tok.Pos.Position("../hello.ugo", code),
		)
	}
}

func loadCode(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}
