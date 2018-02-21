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

// main function.
func main() {
	// variables.
	var subject string

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

	// print log message.
	log.Printf("Subscribing to subject '%s'\n", subject)

	// subscribe messages.
	natsConnection.Subscribe(subject, func(message *nats.Msg) {
		// print log message.
		log.Printf("Received message '%s\n", string(message.Data)+"'")
	})

	// wait for interruption.
	runtime.Goexit()
}
