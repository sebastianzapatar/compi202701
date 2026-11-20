// State: un objeto cambia su comportamiento cuando cambia su estado interno.
// Ejemplo: un escáner con modos Normal / DentroDeString / DentroDeComentario.
package main

import "fmt"

type State interface {
	Read(s *Scanner, ch byte)
}

type Scanner struct {
	state  State
	tokens []string
	buf    []byte
}

type normal struct{}
type inString struct{}
type inComment struct{}

func (normal) Read(s *Scanner, ch byte) {
	switch ch {
	case '"':
		s.state = inString{}
	case '#':
		s.state = inComment{}
	case ' ', '\n':
		// ignorar espacios
	default:
		s.tokens = append(s.tokens, string(ch))
	}
}

func (inString) Read(s *Scanner, ch byte) {
	if ch == '"' {
		s.tokens = append(s.tokens, "STRING("+string(s.buf)+")")
		s.buf = nil
		s.state = normal{}
		return
	}
	s.buf = append(s.buf, ch)
}

func (inComment) Read(s *Scanner, ch byte) {
	if ch == '\n' {
		s.state = normal{}
	}
}

func main() {
	s := &Scanner{state: normal{}}
	src := "x+\"hola mundo\" # esto es un comentario\ny"
	for i := 0; i < len(src); i++ {
		s.state.Read(s, src[i]) // el comportamiento depende del estado actual
	}
	fmt.Println(s.tokens)
}
