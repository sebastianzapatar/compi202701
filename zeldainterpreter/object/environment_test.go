package object

import "testing"

// TestEnclosedEnvironment muestra cómo funcionan los entornos encadenados:
// Get busca hacia afuera, Set escribe solo en el entorno actual.
func TestEnclosedEnvironment(t *testing.T) {
	global := NewEnvironment()
	global.Set("x", &Integer{Value: 1})

	inner := NewEnclosedEnvironment(global)
	inner.Set("y", &Integer{Value: 2})
	inner.Set("x", &Integer{Value: 99}) // "tapa" (shadowing) a la x global

	if v, _ := inner.Get("x"); v.Inspect() != "99" {
		t.Errorf("inner x = %s, se esperaba 99", v.Inspect())
	}
	if v, _ := global.Get("x"); v.Inspect() != "1" {
		t.Errorf("global x = %s, se esperaba 1 (Set no debe tocar outer)", v.Inspect())
	}
	if _, ok := global.Get("y"); ok {
		t.Error("global no debería ver la y del entorno interno")
	}
}

func TestHashKeyByValue(t *testing.T) {
	a, b := &String{Value: "nombre"}, &String{Value: "nombre"}
	if a == b {
		t.Fatal("deben ser punteros distintos para que la prueba tenga sentido")
	}
	if a.HashKey() != b.HashKey() {
		t.Error("dos strings iguales deben tener la misma HashKey")
	}
	if (&Integer{Value: 1}).HashKey() == TrueKey() {
		t.Error("el entero 1 y el booleano true no deben chocar")
	}
}

func TrueKey() HashKey { return (&Boolean{Value: true}).HashKey() }
