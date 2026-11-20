// Goroutines, WaitGroup y Mutex: contar tokens de varios programas en paralelo.
package main

import (
	"fmt"
	"strings"
	"sync"
)

func countTokens(src string) int { return len(strings.Fields(src)) }

func main() {
	programs := []string{
		"let x = 5 ;",
		"let add = fn ( a , b ) { a + b } ;",
		"add ( 1 , 2 ) ;",
	}

	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		total int
	)
	for i, p := range programs {
		wg.Add(1)
		go func() { // desde Go 1.22 cada iteración tiene sus propias i y p
			defer wg.Done()
			n := countTokens(p)
			mu.Lock() // sección crítica: una sola goroutine a la vez
			total += n
			mu.Unlock()
			fmt.Printf("programa %d: %d tokens\n", i, n)
		}()
	}
	wg.Wait() // esperar a que todas terminen
	fmt.Println("total:", total)
}
