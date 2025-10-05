#include <iostream>
#include <thread>
#include <cmath>
#include <iomanip>

using namespace std;

// Interfaz funcional
class FuncionUnivariable {
public:
    virtual double evaluar(double x) const = 0;
    virtual ~FuncionUnivariable() = default;
};

// f(x) = 2x² + 3x + 0.5
class FuncionEjemplo : public FuncionUnivariable {
public:
    double evaluar(double x) const override {
        return 2 * x * x + 3 * x + 0.5;
    }
};

// Clase TrabajadorTrapecio 
class TrabajadorTrapecio {
private:
    const FuncionUnivariable* f;
    double a, h;
    int inicio, fin;
    double* sumaParcial;

public:
    TrabajadorTrapecio(const FuncionUnivariable* f, double a, double h, int inicio, int fin)
        : f(f), a(a), h(h), inicio(inicio), fin(fin) {
        sumaParcial = new double(0.0);
    }

    ~TrabajadorTrapecio() {
        delete sumaParcial;
    }

    void operator()() {
        for (int i = inicio; i <= fin; ++i) {
            double xi = a + i * h;
            *sumaParcial += f->evaluar(xi);
        }
    }

    double getSumaParcial() const { return *sumaParcial; }
};

// Clase Trapecio 
class Trapecio {
public:
    static double integrar(const FuncionUnivariable* f, double a, double b, int n, int numHilos) {
        double h = (b - a) / n;

        // Memoria dinámica para arrays de punteros
        TrabajadorTrapecio** trabajadores = new TrabajadorTrapecio*[numHilos];
        thread** hilos = new thread*[numHilos];

        int tamañoBloque = n / numHilos;

        // Crear objetos trabajadores dinámicamente
        for (int t = 0; t < numHilos; ++t) {
            int inicio = t * tamañoBloque + 1;
            int fin = (t == numHilos - 1) ? n - 1 : (inicio + tamañoBloque - 1);
            trabajadores[t] = new TrabajadorTrapecio(f, a, h, inicio, fin);
        }

        // Lanzar hilos dinámicamente
        for (int t = 0; t < numHilos; ++t)
            hilos[t] = new thread(ref(*trabajadores[t]));

        // Esperar hilos
        for (int t = 0; t < numHilos; ++t) {
            hilos[t]->join();
            delete hilos[t];
        }

        // Calcular suma total
        double suma = f->evaluar(a) + f->evaluar(b);
        for (int t = 0; t < numHilos; ++t) {
            suma += 2 * trabajadores[t]->getSumaParcial();
            delete trabajadores[t];
        }

        // Liberar arrays
        delete[] trabajadores;
        delete[] hilos;

        return (h / 2.0) * suma;
    }
};

int main() {
    FuncionEjemplo f;
    double a = 2.0, b = 20.0;
    int numHilos = 4;

    double anterior = 0.0, actual = 0.0;
    double tolerancia = 1e-9;
    int incremento = 50;
    int n = 1;

    cout << fixed << setprecision(12);

    while (true) {
        actual = Trapecio::integrar(&f, a, b, n, numHilos);
        cout << "n = " << setw(6) << n << "   Área aproximada = " << actual << endl;

        if (n > 1 && fabs(actual - anterior) < tolerancia) {
            cout << "\nEl valor de la integral se ha estabilizado." << endl;
            cout << "Valor final aproximado: " << actual << endl;
            break;
        }

        anterior = actual;
        n += incremento;
    }

    double exacto = 5931.0;
    cout << "\nValor exacto (analítico): " << exacto << endl;
    cout << "Error absoluto: " << fabs(actual - exacto) << endl;

    return 0;
}
