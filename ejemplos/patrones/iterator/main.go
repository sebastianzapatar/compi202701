// Iterator: recorrer una colección sin exponer su estructura interna.
// El lexer ya es un iterador (NextToken). Con Go 1.23 también podemos usar range sobre funciones.
package main

import (
	"fmt"
	"iter"
	"strings"
)

type Token struct{ Type, Literal string }

type Lexer struct {
	words []string
	pos   int
}

func New(src string) *Lexer { return &Lexer{words: strings.Fields(src)} }

// Estilo clásico del libro: llamar NextToken() hasta EOF.
func (l *Lexer) NextToken() Token {
	if l.pos >= len(l.words) {
		return Token{"EOF", ""}
	}
	w := l.words[l.pos]
	l.pos++
	if w == "let" {
		return Token{"LET", w}
	}
	return Token{"WORD", w}
}

// Estilo Go 1.23: un iter.Seq que funciona con for range.
func (l *Lexer) All() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		for {
			t := l.NextToken()
			if t.Type == "EOF" || !yield(t) {
				return
			}
		}
	}
}

func main() {
	l := New("let x = 5 ;")
	for tok := l.NextToken(); tok.Type != "EOF"; tok = l.NextToken() {
		fmt.Print(tok.Type, " ")
	}
	fmt.Println()

	for tok := range New("let y = x ;").All() {
		fmt.Print(tok.Literal, " ")
		if tok.Literal == "=" {
			break // el iterador se detiene limpiamente
		}
	}
	fmt.Println()
}
