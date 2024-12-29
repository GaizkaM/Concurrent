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

// Delator: al llegar informa sobre la llegada de la policía y se va
func main() {

	// Conexión con RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

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

	// Mensaje del delator
	log.Printf("No sóm fumador. ALERTA! Que ve la policia!")

	bodyDelator := "pD"

	err = chPeticionsClients.Publish(
		"",     // Exchange
		s.Name, // Routing key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(bodyDelator),
		})
	failOnError(err, "Failed to publish a message (chPeticionsClients)")

	// Simulación de retardo
	time.Sleep(3 * time.Second)

	for i := 0; i < 3; i++ {
		log.Printf("·")
		time.Sleep(1 * time.Second)
	}
}
