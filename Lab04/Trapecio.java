@FunctionalInterface
interface FuncionUnivariable {
    double evaluar(double x);
}

class TrabajadorTrapecio extends Thread {
    private final FuncionUnivariable f;
    private final double a, h;
    private final int inicio, fin;
    private double sumaParcial = 0;

    public TrabajadorTrapecio(FuncionUnivariable f, double a, double h, int inicio, int fin) {
        this.f = f;
        this.a = a;
        this.h = h;
        this.inicio = inicio;
        this.fin = fin;
    }

    @Override
    public void run() {
        double suma = 0;
        for (int i = inicio; i <= fin; i++) {
            suma += f.evaluar(a + i * h);
        }
        sumaParcial = suma;
    }

    public double getSumaParcial() {
        return sumaParcial;
    }
}

public class Trapecio {
    // Método que calcula la integral con el método del trapecio paralelo
    public static double integrar(FuncionUnivariable f, double a, double b, int n, int numHilos)
            throws InterruptedException {
        double h = (b - a) / n;
        TrabajadorTrapecio[] hilos = new TrabajadorTrapecio[numHilos];
        int tamañoBloque = n / numHilos;

        for (int t = 0; t < numHilos; t++) {
            int inicio = t * tamañoBloque + 1; // evitar f(a)
            int fin = (t == numHilos - 1) ? n - 1 : (inicio + tamañoBloque - 1); // evitar f(b)
            hilos[t] = new TrabajadorTrapecio(f, a, h, inicio, fin);
            hilos[t].start();
        }

        double suma = f.evaluar(a) + f.evaluar(b);
        for (TrabajadorTrapecio hilo : hilos) {
            hilo.join();
            suma += 2 * hilo.getSumaParcial();
        }

        return (h / 2) * suma;
    }
    public static void main(String[] args)  {}
}
