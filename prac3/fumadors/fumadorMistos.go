package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

// Manejo de errores
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// fumadorMistos: solicita mistos al Estanquer y los usa
// Finaliza con el mensaje del Delator
func main() {

	// Conexión con RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Creación de un canal para recibir mistos
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

	// Limitación de un mensaje a la vez
	err = chMistos.Qos(1, 0, false)
	failOnError(err, "Failed to set QoS (chMistos)")

	// Consumición de mensajes de la cola de mistos
	msgsM, err := chMistos.Consume(
		r.Name, // Nombre de la cola
		"",     // Consumer
		false,  // Auto-ack
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Args
	)
	failOnError(err, "Failed to register a consumer (chMistos)")

	// PETICIONES

	// Creación de un canal para enviar peticiones
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

	// Variables
	salir := false

	// ---------------------- INICIO DEL PROCESO ----------------------

	log.Printf("Sóm fumador. Tinc tabac però me falten mistos")

	// Canal para mantener el proceso activo
	forever := make(chan bool)

	for {

		// Publicación de una petición de mistos al estanquer
		bodyDMistos := "dM"
		err = chPeticionsClients.Publish(
			"",     // Exchange
			s.Name, // Routing key
			false,  // Mandatory
			false,  // Immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(bodyDMistos),
			})
		failOnError(err, "Failed to publish a message (chPeticionsClients)")

		// Espera a recibir el Misto
		for m := range msgsM {

			// Si el mensaje recibido NO es un aviso de policía
			if string(m.Body) != "policia" {

				// Reconoce el mensaje
				m.Ack(false)

				// Mensaje informativo de que ha recibido un misto
				log.Printf("He agafat el misto %s. Gràcies!", m.Body)

				// Simulación de retardo
				for i := 0; i < 3; i++ {
					log.Printf("·")
					time.Sleep(1 * time.Second)
				}

				// Solicitud de otro misto
				log.Printf("Me dones un altre misto?")

			} else {

				// Si el mensaje es un aviso de policía
				salir = true

				// Simulación de retardo
				time.Sleep(3 * time.Second)

				// Mensaje de salida
				log.Printf("Anem que ve la policia!")
			}
			break
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
