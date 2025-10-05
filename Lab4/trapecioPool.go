package main

import (
	"fmt"      // Para impresión en consola
	"runtime"  // Para detectar el número de núcleos disponibles
	"sync"     // Para sincronización con WaitGroup y Mutex
	"time"     // Para medir tiempos de ejecución
)

// ===============================================
// DEFINICIÓN DE UN THREAD POOL PERSONALIZADO
// ===============================================
type ThreadPool struct {
	jobs    chan func()     // Canal para recibir funciones a ejecutar
	wg      sync.WaitGroup  // WaitGroup para esperar que terminen los workers
	workers int             // Número de trabajadores (hilos)
}

// Constructor del ThreadPool
func NewThreadPool(workers int) *ThreadPool {
	pool := &ThreadPool{
		jobs:    make(chan func(), workers*10), // Canal con buffer
		workers: workers,
	}

	// Crear los workers como goroutines que esperan funciones del canal
	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done() // Marca que este worker ha terminado cuando sale
			for job := range pool.jobs {
				job() // Ejecuta la función recibida
			}
		}()
	}

	return pool
}

// Método para enviar un trabajo (función) al pool
func (p *ThreadPool) Submit(job func()) {
	p.jobs <- job
}

// Espera a que todos los workers terminen y cierra el canal
func (p *ThreadPool) Wait() {
	close(p.jobs)   // Cierra el canal de trabajos (los workers saldrán del bucle)
	p.wg.Wait()     // Espera a que todos los workers terminen
}

// ===============================================
// ESTRUCTURA PARA CALCULAR INTEGRALES POR TRAPECIO
// ===============================================
type TrapecioPool struct {
	funcion func(float64) float64 // Función a integrar
	a       float64               // Límite inferior
	b       float64               // Límite superior
}

// Constructor para TrapecioPool
func NewTrapecioPool(f func(float64) float64, a, b float64) *TrapecioPool {
	return &TrapecioPool{
		funcion: f,
		a:       a,
		b:       b,
	}
}

// ===============================================
// MÉTODO PARA CALCULAR EL ÁREA CON THREAD POOL
// ===============================================
func (tp *TrapecioPool) CalcularConPool(n int, pool *ThreadPool) float64 {
	h := (tp.b - tp.a) / float64(n)               // Tamaño de cada subintervalo
	numTareas := runtime.NumCPU() * 4             // Más tareas que núcleos para mejorar carga
	var wg sync.WaitGroup                         // WaitGroup para sincronizar las tareas
	resultados := make(chan float64, numTareas)   // Canal para recibir sumas parciales
	suma := 0.5 * (tp.funcion(tp.a) + tp.funcion(tp.b)) // Suma inicial con los extremos

	puntosPorTarea := n / numTareas // Cuántos puntos calcula cada tarea

	// Crear tareas y enviarlas al thread pool
	for tareaID := 0; tareaID < numTareas; tareaID++ {
		wg.Add(1)
		start := tareaID * puntosPorTarea + 1
		end := (tareaID + 1) * puntosPorTarea
		if tareaID == numTareas-1 {
			end = n - 1 // La última tarea puede tomar más puntos
		}

		startLocal := start
		endLocal := end

		// Enviamos una tarea al pool
		pool.Submit(func() {
			defer wg.Done()  // Marca tarea como completada
			sumaLocal := 0.0
			for i := startLocal; i <= endLocal; i++ {
				x := tp.a + float64(i)*h
				sumaLocal += tp.funcion(x)
			}
			resultados <- sumaLocal // Enviamos el resultado parcial al canal
		})
	}

	// Goroutine para cerrar el canal una vez terminadas todas las tareas
	go func() {
		wg.Wait()
		close(resultados)
	}()

	// Acumulamos todos los resultados recibidos desde las tareas
	for resultado := range resultados {
		suma += resultado
	}

	return suma * h // Multiplicamos por el tamaño de subintervalo
}

// ===============================================
// FUNCIÓN PRINCIPAL
// ===============================================
func main() {
	// Definimos la función a integrar
	funcion := func(x float64) float64 {
		return 2*x*x + 3*x + 0.5
	}

	// Detectamos cantidad de núcleos disponibles
	numCores := runtime.NumCPU()

	// Creamos el thread pool
	pool := NewThreadPool(numCores)
	defer pool.Wait() // Esperamos que todos los trabajos terminen al final

	// Creamos la calculadora de integrales
	calculadora := NewTrapecioPool(funcion, 2, 20)

	// Información básica
	fmt.Printf("Usando Thread Pool con %d threads\n", numCores)
	fmt.Println("==============================================")

	// Ejecutamos para diferentes cantidades de divisiones (n)
	for n := 10; n <= 1000000; n *= 10 {
		start := time.Now() // Tiempo de inicio
		area := calculadora.CalcularConPool(n, pool)
		tiempo := time.Since(start) // Tiempo transcurrido

		// Mostramos resultado y tiempo en milisegundos
		fmt.Printf("n=%d: Área=%.6f, Tiempo=%.3f ms\n", 
			n, area, float64(tiempo.Nanoseconds())/1e6)
	}
}
