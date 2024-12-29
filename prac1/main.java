// Projecte creat per Gaizka Medina Gordo

import java.util.Random;
import java.util.concurrent.Semaphore;

public class P1 {

    // Codis ANSI per a donar color als missatges 
    public static final String ANSI_RED = "\u001B[31m";
    public static final String ANSI_GREEN = "\u001B[32m";
    public static final String ANSI_PURPLE = "\u001B[35m";
    public static final String ANSI_CYAN = "\u001B[36m";

    // Arrays de noms d'elfs i rens
    private static String[] NOM_ELFS = {"Taleasin", "Halafarin", "Alduin", "Adamar", "Galather", "Estelar"};
    private static String[] NOM_RENS = {"RUDOLPH", "BLITZEN", "DONDER", "CUPID", "COMET", "VIXEN", "PRANCER", "DANCER", "DASHER"};

    // Constants
    private static final int NUM_ELFS = 6;
    private static final int NUM_RENS = 9;
    private static final int NUM_JOGUINES = 3;

    // Variables
    private static int numElfsEspera = 0; // Número d'elfs a la sala d'espera
    private static int numJoguinesFetes = 0; // Número de joguines fetes en total
    private static int numRensArribats = 0; // Número de rens que han arribat
    private static int numRensLlestos = 0; // Número de rens llestos per partir

    // Semàfors per a la sincronització
    // Semàfor per a la d'espera per a 3 elfs
    private static final Semaphore semEspera = new Semaphore(NUM_JOGUINES);
    // Semàfor per a cada un dels elfs dins la sala de espera 
    private static final Semaphore semElfEsperant = new Semaphore(0);
    // Semàfor per a bloquejar al Pare Noel fins que no hi hagin 3 elfs
    private static final Semaphore semConsultaElfs = new Semaphore(0);
    
    // Semàfor per a cada un dels rens una vegada an acaabt de pasturar
    private static final Semaphore semRenEsperant = new Semaphore(0);
    // Semàfor per a bloquetjar al Pare Noel fins que no estiguin llestos els rens
    private static final Semaphore semConsultaRen = new Semaphore(0);

    // Mutex per a garantir l'exclusió mútua dins la secció crítica de la sala d'espera
    private static final Semaphore mutexElfs = new Semaphore(1, true); 

    public static void main(String[] args) {
        System.out.println(ANSI_CYAN+"SIMULACIÓ DEL PARE NOEL I ELS ELFS EN PRÀCTIQUES:");
        // Crear i iniciar fils per als elfs
        for (int i = 0; i < NUM_ELFS; i++) {
            new Thread(new Elf(NOM_ELFS[i])).start();
        }

        // Crear i iniciar fils per als rens
        for (int i = 0; i < NUM_RENS; i++) {
            new Thread(new Ren(NOM_RENS[i])).start();
        }

        // Fil del Pare Noel
        new Thread(new PareNoel()).start();

    }

    // Classe Elf que implementa Runnable on es duen a terme tots els mètodes i missatges relacionats amb els elfs 
    static class Elf implements Runnable {

        private Random random = new Random();

        private String nom;

        public Elf(String nom) {
            this.nom = nom;
        }

        @Override
        public void run() {
            try {
                // Missatge de començament de la simulació d'un elf
                System.out.println(ANSI_GREEN + "Hola som l'elf " + nom + " i construiré " + NUM_JOGUINES + " joguines.");

                // Bucle que es repiteix el nombre de vegades de les joguines que ha de construir un elf
                for (int i = 1; i < NUM_JOGUINES + 1; i++) {

                    // Simulació de temps de treball de un elf
                    int temps = random.nextInt(100);
                    Thread.sleep(temps);

                    // Semàfor que bloqueja la sala d'espera
                    semEspera.acquire();
                    System.out.println(ANSI_GREEN + nom + " diu: tinc dubtes amb la joguina " + i + ".");

                    // Mutex que protegeix la secció crítica de la sala d'espera (volem exclusió mútua entre els elfs)
                    mutexElfs.acquire();
                    // Augmenta el nombre d'elfs a la sala d'espera
                    numElfsEspera++;
                    // Verificació del nombre d'elfs
                    if (numElfsEspera == 3) {
                        // Si és el tercer elf, mostra un missatge i desperta al Pare Noel
                        System.out.println(ANSI_GREEN + nom + " diu: Som 3 que tenim dubtes, PARE NOEEEEEL!");
                        semConsultaElfs.release();
                    }
                    // Sortida de la secció crítica
                    mutexElfs.release();

                    // Es torna a bloquetjar la sala d'espera
                    semElfEsperant.acquire();

                    // Mutex que protegeix la secció crítica de la sala d'espera (volem exclusió mútua entre els elfs)
                    mutexElfs.acquire();
                    // L'elf construeix la joguina
                    System.out.println(ANSI_GREEN + nom + " diu: Construeixo la joguina amb ajuda.");
                    numJoguinesFetes++;
                    numElfsEspera--;
                    mutexElfs.release();
                    // Es verifica si és l'últim elf dins la sala d'espera
                    if (numElfsEspera == 0) {
                        semEspera.release(3);
                    }

                }
                // L'elf ha finalitzat les seves joguines
                System.out.println(ANSI_GREEN + "L'elf " + nom + " ha fet les seves joguines i acaba <---------");

                // Es verifica si ha estat el darrer elf
                if (numJoguinesFetes == (NUM_ELFS * NUM_JOGUINES)) {
                    System.out.println(ANSI_GREEN + nom + " diu: Som el darrer avisaré al Pare Noel");
                    semConsultaElfs.release();
                }

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }
    }

    // Classe Ren que implementa Runnable on es duen a terme tots els mètodes i missatges relacionats amb els rens  
    static class Ren implements Runnable {

        private Random random = new Random();

        private String nom;

        public Ren(String nom) {
            this.nom = nom;
        }

        @Override
        public void run() {
            try {
                // Missatge de començament de la simulació d'un elf
                System.out.println(ANSI_PURPLE + "El ren " + nom + " se'n va a pasturar.");
                
                // Simula el temps que triga a pasturar
                int temps = random.nextInt(200);
                Thread.sleep(temps);

                // Incrementa el nombre de rens arribats
                numRensArribats++;

                // Verifica si el ren que ha arribat és el darrer
                if (numRensArribats < NUM_RENS) {
                    System.out.println(ANSI_PURPLE + "El ren " + nom + " arriba, " + numRensArribats);
                } else {
                    // Si ho és, mostra el missatge corresponent i desperta al Pare Noel
                    System.out.println(ANSI_PURPLE + "El ren " + nom + " diu: Som el darrer en voler, podem partir.");
                    semConsultaRen.release(); 
                }

                // Els rens esperen a ser enganxats
                semRenEsperant.acquire(); 
                numRensLlestos++;
                System.out.println(ANSI_PURPLE + "El ren " + nom + " està enganxat al trineu.");

                // Verifica si el ren enganxat és el darrer
                if (numRensLlestos == NUM_RENS) {
                    // Avisa al Pare Noel
                    semConsultaRen.release();
                } else {
                    // Si no és el darrer, s'allibera el ren
                    semRenEsperant.release();
                }

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }
    }

    // Classe Pare Noel que implementa Runnable on es duen a terme tots els mètodes i missatges relacionats amb el Pare Noel
    static class PareNoel implements Runnable {

        @Override
        public void run() {
            try {
                // Missatge de començament de la simulació del Pare Noel
                System.out.println(ANSI_RED + "-------> El Pare Noel diu: Estic despert però me'n torn a jeure.");

                // Bloqueig del Pare Noel. No continuarà fins que no hi hagin 3 elfs amb dubtes
                semConsultaElfs.acquire();

                // Bucle while fins que s'hagin fet totes les joguines
                while (numJoguinesFetes < (NUM_ELFS * NUM_JOGUINES)) {
                    
                    // Si hi ha 3 elfs a la sala d'espera, mostra un missatge i desperta als 3 elfs
                    System.out.println(ANSI_RED + "-------> El Pare Noel diu: Atendré els dubtes d'aquests 3.");
                    semElfEsperant.release(3);

                    // Es torna a dormir
                    System.out.println(ANSI_RED + "-------> El Pare Noel diu: Estic cansat me'n torn a jeure");
                    semConsultaElfs.acquire();
                }

                // Una vegada han acabat les joguines, verifica els rens
                System.out.println(ANSI_RED + "-------> Pare Noel diu: Les joguines estan llestes. I Els rens?");
                // Bloqueig fins que estiguin els rens
                semConsultaRen.acquire();

                // Una vegada estan els rens, s'han d'enganxar
                System.out.println(ANSI_RED + "-------> El Pare Noel diu: Enganxaré els rens i partiré.");
                // S'allibera el primer ren
                semRenEsperant.release();
                // Es bloqueja fins s'alliberen (enganxen) tots els rens
                semConsultaRen.acquire();

                // Missatge de fi de la simulació
                System.out.println(ANSI_RED + "-------> El Pare Noel ha enganxat els rens, ha carregat les joguines i se'n va.");

            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }
    }
}
