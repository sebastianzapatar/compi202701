// Worker pool + fan-out/fan-in + context: evaluar muchos programas en paralelo con timeout.
package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Job struct {
	ID  int
	Src string
}

type Result struct {
	ID    int
	Value int
	Err   error
}

// evalSum "evalúa" programas de la forma "1 + 2 + 3" y simula trabajo lento.
func evalSum(ctx context.Context, src string) (int, error) {
	select {
	case <-time.After(time.Duration(len(src)) * 5 * time.Millisecond):
	case <-ctx.Done():
		return 0, ctx.Err() // cancelado o se acabó el tiempo
	}
	total := 0
	for _, part := range strings.Split(src, "+") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return 0, fmt.Errorf("token inválido %q", strings.TrimSpace(part))
		}
		total += n
	}
	return total, nil
}

func worker(ctx context.Context, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		v, err := evalSum(ctx, j.Src)
		results <- Result{ID: j.ID, Value: v, Err: err}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	programs := []string{"1 + 2", "10 + 20 + 30", "4 + x", "1 + 1 + 1 + 1 + 1 + 1 + 1 + 1 + 1"}
	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup
	for range 2 { // fan-out: 2 workers comparten el canal jobs
		wg.Add(1)
		go worker(ctx, jobs, results, &wg)
	}

	go func() { // productor
		defer close(jobs)
		for i, p := range programs {
			jobs <- Job{ID: i, Src: p}
		}
	}()

	go func() { // fan-in: cerrar results cuando todos los workers terminen
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.Err != nil {
			fmt.Printf("programa %d: error: %v\n", r.ID, r.Err)
			continue
		}
		fmt.Printf("programa %d = %d\n", r.ID, r.Value)
	}
}
