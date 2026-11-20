// Observer: un objeto notifica a sus suscriptores cuando ocurre algo.
// Ejemplo: el evaluador emite eventos; un tracer y un profiler los escuchan.
package main

import "fmt"

type Event struct {
	Kind string // "call" | "error"
	Name string
}

type Observer interface{ OnEvent(e Event) }

type Evaluator struct{ observers []Observer }

func (ev *Evaluator) Subscribe(o Observer) { ev.observers = append(ev.observers, o) }

func (ev *Evaluator) emit(e Event) {
	for _, o := range ev.observers {
		o.OnEvent(e)
	}
}

func (ev *Evaluator) Call(fn string) {
	ev.emit(Event{Kind: "call", Name: fn})
	if fn == "foo" {
		ev.emit(Event{Kind: "error", Name: "identifier not found: foo"})
	}
}

type Tracer struct{}

func (Tracer) OnEvent(e Event) { fmt.Printf("[trace] %-5s %s\n", e.Kind, e.Name) }

type Profiler struct{ calls map[string]int }

func (p *Profiler) OnEvent(e Event) {
	if e.Kind == "call" {
		p.calls[e.Name]++
	}
}

func main() {
	ev := &Evaluator{}
	prof := &Profiler{calls: map[string]int{}}
	ev.Subscribe(Tracer{})
	ev.Subscribe(prof)

	for _, f := range []string{"fib", "fib", "len", "foo"} {
		ev.Call(f)
	}
	fmt.Println("perfil:", prof.calls)
}
