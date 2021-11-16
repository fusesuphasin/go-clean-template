package infrastructure

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// Here we set the way error messages are displayed in the terminal.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func RabbitConn() (ch *amqp.Connection, err error) {
	connRabbitMQ, err := amqp.Dial("amqp://" + os.Getenv("RABBITMQ_USER") + ":" + os.Getenv("RABBITMQ_PASS") + "@" + os.Getenv("RABBITMQ_HOST") + ":" + os.Getenv("RABBITMQ_PORT") + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer connRabbitMQ.Close()

	return connRabbitMQ, err

	// return ch, errCh
}

func SendQueue(data string, queuename string) {
	fmt.Println("data : ", data)
	//str := fmt.Sprintf("%v", data)
	conn, _ := RabbitConn()
	defer conn.Close()
	ch, errCh := conn.Channel()
	if errCh != nil {
		panic(errCh)
	}
	_, err := ch.QueueDeclare(
		queuename,		// name
		true,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	failOnError(err, "Failed to declare a queue")

	// Attempt to publish a message to the queue.
	err = ch.Publish(
		"",
		queuename,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		},
	)
	if err != nil {
		panic(err)
	}

	defer ch.Close()
}

func GetQueue(queuename string) /* string */{
	//var message string
	conn, _ := RabbitConn()
	defer conn.Close()
	ch, errCh := conn.Channel()
	if errCh != nil {
		panic(errCh)
	}

	q, err := ch.QueueDeclare(
		queuename,		// name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
	failOnError(err, "Failed to register a consumer")
	
	//forever := make(chan bool) 
	go func() {
		for d := range msgs {
		  fmt.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("_____________________")
	//<-forever
	//return message
}