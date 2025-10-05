#include <iostream>      // Para entrada y salida estándar
#include <thread>        // Para manejar múltiples hilos
#include <vector>        // Para usar el contenedor vector
#include <chrono>        // Para medir tiempos de ejecución
#include <cmath>         // Para funciones matemáticas
#include <mutex>         // Para evitar condiciones de carrera con múltiples hilos

using namespace std;

// -------------------------------------
// FUNCION A INTEGRAR: f(x) = 2x² + 3x + 0.5
// -------------------------------------
double funcion(double x) {
    return 2 * x * x + 3 * x + 0.5;
}

// -------------------------------------------------
// FUNCION SECUENCIAL: Método del trapecio (versión no paralela)
// -------------------------------------------------
double calcularAreaSecuencial(double a, double b, int n) {
    double h = (b - a) / n; // Tamaño del subintervalo
    double suma = 0.5 * (funcion(a) + funcion(b)); // Sumamos extremos (regla del trapecio)
    
    // Sumamos las áreas internas
    for (int i = 1; i < n; i++) {
        double x = a + i * h;
        suma += funcion(x);
    }
    
    return suma * h; // Multiplicamos la suma por el tamaño del paso (área total)
}

// -------------------------------------------------
// FUNCION PARALELA: Método del trapecio usando múltiples hilos
// -------------------------------------------------
double calcularAreaParalela(double a, double b, int n, int numThreads) {
    double h = (b - a) / n; // Tamaño del subintervalo
    double sumaTotal = 0.5 * (funcion(a) + funcion(b)); // Sumamos extremos
    vector<thread> threads; // Vector para guardar los hilos creados
    mutex mtx; // Mutex para proteger acceso a sumaTotal entre hilos
    
    int puntosPorThread = n / numThreads; // Cuántos puntos calcula cada hilo

    // Creamos hilos y asignamos rangos de cálculo a cada uno
    for (int t = 0; t < numThreads; t++) {
        threads.emplace_back([&, t]() { // Captura por referencia (usar "&")
            int start = t * puntosPorThread + 1;
            int end = (t == numThreads - 1) ? n - 1 : (t + 1) * puntosPorThread;
            
            double sumaLocal = 0;
            // Cada hilo calcula su suma parcial
            for (int i = start; i <= end; i++) {
                double x = a + i * h;
                sumaLocal += funcion(x);
            }

            // Sección crítica: sumamos la parte local al total global (protegido con mutex)
            lock_guard<mutex> lock(mtx);
            sumaTotal += sumaLocal;
        });
    }

    // Esperamos que todos los hilos terminen
    for (auto& thread : threads) {
        thread.join();
    }

    return sumaTotal * h; // Multiplicamos suma total por h (resultado final)
}

// -------------------------------------------------
// FUNCIÓN PRINCIPAL
// -------------------------------------------------
int main() {
    double a = 2.0, b = 20.0; // Intervalo de integración
    int numCores = thread::hardware_concurrency(); // Detectamos cantidad de núcleos disponibles

    if (numCores == 0) numCores = 2; // Si no se puede detectar, usamos 2 como valor por defecto

    // Información básica
    cout << "Cálculo de integral usando " << numCores << " threads" << endl;
    cout << "Función: 2x² + 3x + 0.5" << endl;
    cout << "Intervalo: [" << a << ", " << b << "]" << endl;
    cout << "=====================================" << endl;

    // Repetimos el cálculo con distintos valores de 'n' (más puntos => más precisión y más trabajo)
    for (int n = 10; n <= 1000000; n *= 10) {
        // -------- Versión secuencial --------
        auto start = chrono::high_resolution_clock::now(); // Tiempo de inicio
        double areaSec = calcularAreaSecuencial(a, b, n);  // Cálculo secuencial
        auto end = chrono::high_resolution_clock::now();   // Tiempo de fin
        auto tiempoSec = chrono::duration_cast<chrono::microseconds>(end - start); // Duración

        // -------- Versión paralela --------
        start = chrono::high_resolution_clock::now(); // Tiempo de inicio
        double areaPar = calcularAreaParalela(a, b, n, numCores); // Cálculo paralelo
        end = chrono::high_resolution_clock::now();   // Tiempo de fin
        auto tiempoPar = chrono::duration_cast<chrono::microseconds>(end - start); // Duración

        // -------- Mostrar resultados --------
        cout << "n=" << n
             << ":\tSecuencial=" << tiempoSec.count()/1000.0 << "ms" 
             << "\tParalelo=" << tiempoPar.count()/1000.0 << "ms"
             << "\tSpeedup=" << tiempoSec.count() / (double)tiempoPar.count() << "x" << endl;
    }

    return 0; // Fin del programa
}
