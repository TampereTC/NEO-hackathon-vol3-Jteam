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
	var subscribeTopic string
	var publishTopic string
	var messageAlley = make(chan []byte, 10)

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

	// set nats subscribe topic name.
	subscribeTopic = "foo"

	// execute message handler.
	go func() {
		// variables.
		var receivedMessage message.Message
		var sentMessage message.Message
		var messageCounter int
		var jsonMessage []byte

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

					// last message received.
					if receivedMessage.LastMessage == true {
						// set nats publish topic name.
						publishTopic = "bar"

						// set sent message.
						messageCounter = receivedMessage.MessageCounter + 1
						sentMessage.MessageCounter = messageCounter
						sentMessage.MessageText = "Legit omnia mandata."
						sentMessage.LastMessage = true

						// marshall message.
				    jsonMessage, _ = json.Marshal(sentMessage)

				    // publish message.
				    natsConnection.Publish(publishTopic, jsonMessage)

						// print log message.
						log.Println("Completion message sent.")
					}
				// default action.
				default:
					// do nothing.
			}
		}
	}()

	// print log message.
	log.Printf("Subscribing from '%s'.\n", subscribeTopic)

	// subscribe messages.
	natsConnection.Subscribe(subscribeTopic, func(message *nats.Msg) {
		// forward message data for processing.
		messageAlley <- message.Data
	})

	// wait for interruption.
	runtime.Goexit()
}
