package main

import (
	"log"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

// Manejo de errores
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Estanquer: espera sin hacer nada hasta que vienen clientes fumadores y les proporciona tabaco/mistos
// Finaliza con el mensaje del Delator
func main() {

	// Conexión con RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Creación de un canal para enviar el tabaco
	chTabac, err := conn.Channel()
	failOnError(err, "Failed to open a channel (chTabac)")
	defer chTabac.Close()

	// Declaración de la cola correspondiente
	q, err := chTabac.QueueDeclare(
		"tQueue", // Nombre de la cola
		false,    // Durable
		true,     // Auto-Delete
		false,    // Exclusive
		false,    // No-wait
		nil,      // Arguments
	)
	failOnError(err, "Failed to declare a queue (chTabac)")

	// Creación de un canal para enviar los mistos
	chMistos, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer chMistos.Close()

	// Declaración de la cola correspondiente
	r, err := chMistos.QueueDeclare(
		"mQueue", // Nombre de la cola
		false,    // Durable
		true,     // Auto-Delete
		false,    // Exclusive
		false,    // No-wait
		nil,      // Arguments
	)
	failOnError(err, "Failed to declare a queue (chMistos)")

	// Creación de un canal para recibir peticiones
	chPeticionsClients, err := conn.Channel()
	failOnError(err, "Failed to open a channel (chPeticionsClients)")
	defer chPeticionsClients.Close()

	// Declaración de la cola correspondiente
	s, err := chPeticionsClients.QueueDeclare(
		"pQueue", // Nombre de la cola
		false,    // Durable
		true,     // Auto-Delete
		false,    // Exclusive
		false,    // No-wait
		nil,      // Arguments
	)
	failOnError(err, "Failed to declare a queue (chPeticionsClients)")

	// Declaración de una cola para poder recoger peticiones
	msgsP, err := chPeticionsClients.Consume(
		s.Name, // Nombre de la cola
		"",     // Consumer
		false,  // Auto-ack
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Args
	)

	// Variables
	indexT := 0
	indexM := 0
	salir := false

	// ---------------------- INICIO DEL PROCESO ----------------------

	log.Printf("Hola, som l'estanquer il·legal")

	// Canal para mantener el proceso activo
	forever := make(chan bool)

	// Espera a tener una petición
	for f := range msgsP {

		// Si la petición es de tabaco
		if string(f.Body) == "dT" {

			// Incrementa el contador de tabaco
			indexT++
			// Reconoce el mensaje
			f.Ack(false)

			bodyT := strconv.Itoa(indexT)

			err = chTabac.Publish(
				"",     // Exchange
				q.Name, // Routing key
				false,  // Mandatory
				false,  // Immediate
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(bodyT),
				})
			failOnError(err, "Failed to publish a message (chTabac)")

			// Mensaje avisando que ha puesto el tabaco
			log.Printf("He posat el tabac %s damunt la taula", bodyT)

			// Si la petición es de mistos
		} else if string(f.Body) == "dM" {

			// Incrementa el contador de mistos
			indexM++
			// Reconoce el mensaje
			f.Ack(false)

			bodyM := strconv.Itoa(indexM)

			err = chMistos.Publish(
				"",     // Exchange
				r.Name, // Routine Key
				false,  // Mandatory
				false,  // Immediate
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(bodyM),
				})
			failOnError(err, "Failed to publish a message (chMistos)")

			// Mensaje avisando que ha puesto el misto
			log.Printf("He posat el misto %s damunt la taula", bodyM)

			// Sino, ha recibido el mensaje del delator
		} else {

			// Reconoce el mensaje
			f.Ack(false)

			// Condición de salida
			salir = true

			body := "policia"

			// Publica el mensaje de aviso en ambas colas
			// Tabaco
			err = chTabac.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(body),
				})
			failOnError(err, "Failed to publish a message (chMistos)")

			// Mistos
			err = chMistos.Publish(
				"",
				r.Name,
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         []byte(body),
				})
			failOnError(err, "Failed to publish a message (chMistos)")

			// Simulación de retardo
			time.Sleep(3 * time.Second)

			// Mensaje de finalización
			log.Printf("Uyuyuy la policia! Men vaig")
			for i := 0; i < 3; i++ {
				log.Printf("·")
				time.Sleep(1 * time.Second)
			}
			log.Printf("Men duc la taula!!!")
		}

		// Si debe salir, rompe el bucle
		if salir {
			break
		}
	}

	// Si no debe salir inmediatamente, mantiene el proceso activo
	if !salir {
		<-forever
	}
}
