package main

import (
	"fmt"
	"math"
	"sync"
)

// Interfaz funcional
type FuncionUnivariable interface {
	Evaluar(x float64) float64
}

// f(x) = 2x² + 3x + 0.5
type FuncionEjemplo struct{}

func (f FuncionEjemplo) Evaluar(x float64) float64 {
	return 2*x*x + 3*x + 0.5
}

// Estructura de trabajador 
type TrabajadorTrapecio struct {
	f           FuncionUnivariable
	a, h        float64
	inicio, fin int
	sumaParcial float64
}

func (t *TrabajadorTrapecio) Calcular(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := t.inicio; i <= t.fin; i++ {
		xi := t.a + float64(i)*t.h
		t.sumaParcial += t.f.Evaluar(xi)
	}
}

// Función de integración 
func integrar(f FuncionUnivariable, a, b float64, n, numGoroutines int) float64 {
	h := (b - a) / float64(n)

	trabajadores := make([]*TrabajadorTrapecio, numGoroutines)
	var wg sync.WaitGroup

	tamañoBloque := n / numGoroutines

	for t := 0; t < numGoroutines; t++ {
		inicio := t*tamañoBloque + 1
		fin := n - 1
		if t != numGoroutines-1 {
			fin = inicio + tamañoBloque - 1
		}
		trabajadores[t] = &TrabajadorTrapecio{f: f, a: a, h: h, inicio: inicio, fin: fin}
		wg.Add(1)
		go trabajadores[t].Calcular(&wg)
	}

	wg.Wait()

	suma := f.Evaluar(a) + f.Evaluar(b)
	for _, trab := range trabajadores {
		suma += 2 * trab.sumaParcial
	}

	return (h / 2) * suma
}

func main() {
	var f FuncionEjemplo
	a, b := 2.0, 20.0
	numGoroutines := 4

	anterior, actual := 0.0, 0.0
	tolerancia := 1e-9
	incremento := 50
	n := 1

	for {
		actual = integrar(f, a, b, n, numGoroutines)
		fmt.Printf("n = %-6d  Área aproximada = %.12f\n", n, actual)

		if n > 1 && math.Abs(actual-anterior) < tolerancia {
			fmt.Println("\nEl valor de la integral se ha estabilizado.")
			fmt.Printf("Valor final aproximado: %.12f\n", actual)
			break
		}

		anterior = actual
		n += incremento
	}

	exacto := 5931.0
	fmt.Printf("\nValor exacto (analítico): %.12f\n", exacto)
	fmt.Printf("Error absoluto: %.12f\n", math.Abs(actual-exacto))
}
