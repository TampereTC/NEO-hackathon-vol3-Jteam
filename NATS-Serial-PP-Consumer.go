package main

/*
  The original application is in https://nats.io/documentation/tutorials/nats-client-dev/,
  and it has been used almost as is. I've extended some variable names and added comments
  to clarify my own understanding.
*/

// import packages.
import "log"
import "os"
import "runtime"
import "github.com/nats-io/go-nats"
import "encoding/json"
import "./message"

// main function.
func main() {
	// variables.
	var subject string
	var messageAlley = make(chan []byte, 1)

	// connect to nats server.
	natsConnection, connectionError := nats.Connect(nats.DefaultURL)

	// handle connection error.
	if connectionError != nil {
		// print log messages.
		log.Printf("Connection to NATS server '%s' failed.\n", nats.DefaultURL)
		log.Printf("'%s'.\n", connectionError.Error())

		// exit with error.
		os.Exit(1)
	}

	// print log message.
	log.Printf("Connected to '%s'.\n", nats.DefaultURL)

	// set subject name.
	subject = "foo"

	// execute message handler.
	go func() {
		// variables.
		var receivedMessage message.Message

		// begin endless loop.
		for {
			// wait for input.
			select {
				// receive message.
				case receivedJson := <- messageAlley:
					// unmarshal json; []byte.
					json.Unmarshal(receivedJson, &receivedMessage)

					// print message content.
					log.Printf("Counter: %d\n", receivedMessage.MessageCounter)
					log.Printf("Text: %s\n", receivedMessage.MessageText)
					log.Printf("Last: %t\n", receivedMessage.LastMessage)
				// default action.
				default:
					// print no message.
					//log.Println("No message.")
			}
		}
	}()

	// print log message.
	log.Printf("Subscribing to subject '%s'\n", subject)

	// subscribe messages.
	natsConnection.Subscribe(subject, func(message *nats.Msg) {
		// print log message.
		log.Printf("Received message '%s\n", string(message.Data)+"'")
		messageAlley <- message.Data
	})

	// wait for interruption.
	runtime.Goexit()
}
