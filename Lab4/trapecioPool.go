package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type ThreadPool struct {
	jobs    chan func()
	wg      sync.WaitGroup
	workers int
}

func NewThreadPool(workers int) *ThreadPool {
	pool := &ThreadPool{
		jobs:    make(chan func(), workers*10),
		workers: workers,
	}

	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			for job := range pool.jobs {
				job()
			}
		}()
	}

	return pool
}

func (p *ThreadPool) Submit(job func()) {
	p.jobs <- job
}

func (p *ThreadPool) Wait() {
	close(p.jobs)
	p.wg.Wait()
}

type TrapecioPool struct {
	funcion func(float64) float64
	a       float64
	b       float64
}

func NewTrapecioPool(f func(float64) float64, a, b float64) *TrapecioPool {
	return &TrapecioPool{
		funcion: f,
		a:       a,
		b:       b,
	}
}

func (tp *TrapecioPool) CalcularConPool(n int, pool *ThreadPool) float64 {
	h := (tp.b - tp.a) / float64(n)
	numTareas := runtime.NumCPU() * 4

	var wg sync.WaitGroup
	resultados := make(chan float64, numTareas)
	suma := 0.5 * (tp.funcion(tp.a) + tp.funcion(tp.b))

	puntosPorTarea := n / numTareas

	for tareaID := 0; tareaID < numTareas; tareaID++ {
		wg.Add(1)
		start := tareaID * puntosPorTarea + 1
		end := (tareaID + 1) * puntosPorTarea
		if tareaID == numTareas-1 {
			end = n - 1
		}

		startLocal := start
		endLocal := end
		pool.Submit(func() {
			defer wg.Done()
			sumaLocal := 0.0
			for i := startLocal; i <= endLocal; i++ {
				x := tp.a + float64(i)*h
				sumaLocal += tp.funcion(x)
			}
			resultados <- sumaLocal
		})
	}

	go func() {
		wg.Wait()
		close(resultados)
	}()

	for resultado := range resultados {
		suma += resultado
	}

	return suma * h
}

func main() {
	funcion := func(x float64) float64 {
		return 2*x*x + 3*x + 0.5
	}

	numCores := runtime.NumCPU()
	pool := NewThreadPool(numCores)
	defer pool.Wait()

	calculadora := NewTrapecioPool(funcion, 2, 20)

	fmt.Printf("Usando Thread Pool con %d threads\n", numCores)
	fmt.Println("==============================================")

	for n := 10; n <= 1000000; n *= 10 {
		start := time.Now()
		area := calculadora.CalcularConPool(n, pool)
		tiempo := time.Since(start)

		fmt.Printf("n=%d: Área=%.6f, Tiempo=%.3f ms\n", 
			n, area, float64(tiempo.Nanoseconds())/1e6)
	}
}