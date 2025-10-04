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
    public static void main(String[] args)  {}
}
