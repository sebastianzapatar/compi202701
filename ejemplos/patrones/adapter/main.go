// Adapter: hacer que una interfaz existente encaje con la que espera nuestro código.
// Ejemplo: el lexer espera un Source; queremos alimentarlo desde cualquier io.Reader.
package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Source es la interfaz que conoce nuestro lexer.
type Source interface {
	NextChar() (byte, bool)
}

// stringSource es la implementación nativa.
type stringSource struct {
	s   string
	pos int
}

func (s *stringSource) NextChar() (byte, bool) {
	if s.pos >= len(s.s) {
		return 0, false
	}
	s.pos++
	return s.s[s.pos-1], true
}

// readerAdapter adapta un io.Reader (archivo, red, stdin) a Source.
type readerAdapter struct{ r *bufio.Reader }

func FromReader(r io.Reader) Source { return &readerAdapter{r: bufio.NewReader(r)} }

func (a *readerAdapter) NextChar() (byte, bool) {
	b, err := a.r.ReadByte()
	if err != nil {
		return 0, false
	}
	return b, true
}

// countLetters es el cliente: solo depende de Source.
func countLetters(src Source) int {
	n := 0
	for ch, ok := src.NextChar(); ok; ch, ok = src.NextChar() {
		if ch >= 'a' && ch <= 'z' {
			n++
		}
	}
	return n
}

func main() {
	fmt.Println(countLetters(&stringSource{s: "let x = 5;"}))
	fmt.Println(countLetters(FromReader(strings.NewReader("fn(a, b)"))))
}
