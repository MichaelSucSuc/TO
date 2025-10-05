import java.util.concurrent.*; // Importa las clases necesarias para la concurrencia y las tareas en segundo plano.
import java.util.function.Function; // Importa la clase Function que se utiliza para definir la función matemática a integrar.

public class TrapecioPool {
    // Definimos el tamaño del pool de hilos (threads) según la cantidad de procesadores disponibles en la máquina.
    private static final int POOL_SIZE = Runtime.getRuntime().availableProcessors();
    // Creamos un thread pool con el tamaño definido anteriormente.
    private static final ExecutorService threadPool = Executors.newFixedThreadPool(POOL_SIZE);
    
    private Function<Double, Double> funcion; // La función que queremos integrar.
    private double a; // Límite inferior de la integral.
    private double b; // Límite superior de la integral.

    // Constructor de la clase que recibe la función y los límites de integración.
    public TrapecioPool(Function<Double, Double> funcion, double a, double b) {
        this.funcion = funcion;
        this.a = a;
        this.b = b;
    }

    // Método que realiza la integración usando el método de los trapecios y un thread pool para paralelizar el trabajo.
    public double calcularConPool(int n) throws InterruptedException, ExecutionException {
        // Calculamos el tamaño de cada subintervalo
        double h = (b - a) / n;
        // La suma inicial es la mitad de los valores de los extremos (puntos a y b)
        double suma = 0.5 * (funcion.apply(a) + funcion.apply(b));
        
        // Definimos cuántas tareas vamos a dividir. En este caso, usamos el pool de hilos con más tareas que hilos,
        // lo que permite mejor distribución de trabajo entre los hilos.
        int tareas = POOL_SIZE * 4; 
        // Cada tarea calculará una cantidad proporcional de puntos.
        int puntosPorTarea = n / tareas;
        
        // Array para almacenar las tareas que se envían al thread pool.
        Future<Double>[] resultados = new Future[tareas];
        
        // Dividimos el trabajo en subintervalos y enviamos tareas al thread pool.
        for (int t = 0; t < tareas; t++) {
            // Cada tarea calcula una parte del intervalo
            final int start = t * puntosPorTarea + 1;
            final int end = (t == tareas - 1) ? n - 1 : (t + 1) * puntosPorTarea;
            
            // Enviamos la tarea al thread pool para su ejecución en paralelo
            resultados[t] = threadPool.submit(() -> {
                double sumaParcial = 0;
                // Cada hilo calcula una suma parcial de la integral sobre un subintervalo.
                for (int i = start; i <= end; i++) {
                    double x = a + i * h;
                    sumaParcial += funcion.apply(x);
                }
                return sumaParcial;
            });
        }
        
        // Esperamos a que todas las tareas terminen y sumamos sus resultados parciales.
        for (Future<Double> resultado : resultados) {
            suma += resultado.get(); // Espera a que la tarea termine y obtiene el resultado.
        }
        
        // Multiplicamos la suma total por el tamaño del intervalo (h) para obtener el valor de la integral.
        return suma * h;
    }
    
    // Método para cerrar el pool de hilos una vez terminado el cálculo.
    public static void shutdownPool() {
        threadPool.shutdown(); // Cierra el thread pool.
    }

    // Método principal para probar la clase y calcular la integral para diferentes valores de n.
    public static void main(String[] args) {
        // Definimos la función a integrar: 2x^2 + 3x + 0.5
        Function<Double, Double> funcion = x -> 2 * Math.pow(x, 2) + 3 * x + 0.5;
        // Creamos una instancia de la clase TrapecioPool con la función y los límites de integración.
        TrapecioPool calculadora = new TrapecioPool(funcion, 2, 20);
        
        try {
            // Imprimimos información sobre la cantidad de hilos en el pool.
            System.out.println("Usando Thread Pool con " + POOL_SIZE + " threads");
            System.out.println("==============================================");
            
            // Iteramos sobre diferentes valores de n para ver cómo cambia el tiempo de ejecución y el resultado de la integral.
            for (int n = 10; n <= 1000000; n *= 10) {
                long startTime = System.nanoTime(); // Marca el tiempo de inicio.
                double area = calculadora.calcularConPool(n); // Calculamos el área usando el método de Trapecio con thread pool.
                long endTime = System.nanoTime(); // Marca el tiempo de finalización.
                
                // Imprimimos el resultado con el tiempo de ejecución.
                System.out.printf("n=%d: Área=%.6f, Tiempo=%.3f ms\n", 
                    n, area, (endTime - startTime) / 1e6);
            }
            
        } catch (Exception e) {
            // Si ocurre un error, lo mostramos.
            e.printStackTrace();
        } finally {
            // Cerramos el thread pool al finalizar el cálculo.
            shutdownPool();
        }
    }
}
