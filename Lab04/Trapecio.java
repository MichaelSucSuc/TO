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
    public static void main(String[] args) throws InterruptedException {
        FuncionUnivariable f = x -> 2 * x * x + 3 * x + 0.5;
        double a = 2, b = 20;
        int numHilos = 4;

        double anterior = 0.0;
        double actual;
        int n = 1;

        double tolerancia = 1e-9;
        int incremento = 50;

        System.out.println("Cálculo del área bajo la curva f(x) = 2x² + 3x + 0.5");
        System.out.println("Intervalo [" + a + ", " + b + "]\n");

        while (true) {
            actual = integrar(f, a, b, n, numHilos);
            System.out.printf("n = %-6d  Área aproximada = %.12f%n", n, actual);

            if (n > 1 && Math.abs(actual - anterior) < tolerancia) {
                System.out.println("\nEl valor de la integral se ha estabilizado");
                System.out.printf("Valor final aproximado: %.12f%n", actual);
                break;
            }

            anterior = actual;
            n += incremento;
        }

        double exacto = (2.0 / 3.0) * Math.pow(b, 3) + (3.0 / 2.0) * b * b + 0.5 * b
                - ((2.0 / 3.0) * Math.pow(a, 3) + (3.0 / 2.0) * a * a + 0.5 * a);

        System.out.printf("\nValor exacto (analítico): %.12f%n", exacto);
        System.out.printf("Error absoluto: %.12f%n", Math.abs(actual - exacto));
    }
}
