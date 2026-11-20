// Un lexer concurrente: produce tokens en una goroutine y los envía por un canal.
// El consumidor (el parser) los recibe con range hasta que el canal se cierra.
package main

import (
	"fmt"
	"unicode"
)

type Token struct{ Type, Literal string }

func lex(input string) <-chan Token {
	out := make(chan Token, 4) // buffer: el lexer puede adelantarse hasta 4 tokens
	go func() {
		defer close(out) // cerrar avisa al consumidor que no hay más tokens
		for i := 0; i < len(input); {
			ch := rune(input[i])
			switch {
			case unicode.IsSpace(ch):
				i++
			case unicode.IsLetter(ch):
				start := i
				for i < len(input) && unicode.IsLetter(rune(input[i])) {
					i++
				}
				out <- Token{"IDENT", input[start:i]}
			case unicode.IsDigit(ch):
				start := i
				for i < len(input) && unicode.IsDigit(rune(input[i])) {
					i++
				}
				out <- Token{"INT", input[start:i]}
			default:
				out <- Token{"SYMBOL", string(ch)}
				i++
			}
		}
		out <- Token{"EOF", ""}
	}()
	return out
}

func main() {
	for tok := range lex("let total = precio * 3;") {
		fmt.Printf("%-6s %q\n", tok.Type, tok.Literal)
	}
}
