// Objetos en Go: struct + métodos + constructor + encapsulación.
package main

import (
	"errors"
	"fmt"
)

// Token es un "objeto valor": datos inmutables con comportamiento.
type Token struct {
	Type    string
	Literal string
}

// String implementa fmt.Stringer: Println lo usa automáticamente.
func (t Token) String() string { return fmt.Sprintf("%s(%q)", t.Type, t.Literal) }

// Cursor encapsula su estado: los campos en minúscula no se ven fuera del paquete.
type Cursor struct {
	input string
	pos   int
}

// NewCursor es el "constructor": valida y entrega el objeto listo para usar.
func NewCursor(input string) (*Cursor, error) {
	if input == "" {
		return nil, errors.New("input vacío")
	}
	return &Cursor{input: input}, nil
}

// Next avanza el cursor (pointer receiver: modifica el estado).
func (c *Cursor) Next() (byte, bool) {
	if c.pos >= len(c.input) {
		return 0, false
	}
	ch := c.input[c.pos]
	c.pos++
	return ch, true
}

// Pos es un "getter": lectura controlada del estado interno.
func (c *Cursor) Pos() int { return c.pos }

func main() {
	t := Token{Type: "INT", Literal: "5"}
	fmt.Println(t) // usa t.String()

	c, err := NewCursor("let")
	if err != nil {
		panic(err)
	}
	for ch, ok := c.Next(); ok; ch, ok = c.Next() {
		fmt.Printf("%c@%d ", ch, c.Pos()-1)
	}
	fmt.Println()

	if _, err := NewCursor(""); err != nil {
		fmt.Println("error:", err)
	}
}
