package main

/*
  The original application is in https://nats.io/documentation/tutorials/nats-client-dev/,
  and it has been used almost as is. I've extended some variable names and added comments
  to clarify my own understanding.
*/

// import packages.
import "fmt"
import "os"
import "github.com/nats-io/go-nats"
import "time"
import "strconv"

// main function.
func main() {
  // variables.
  var subject string
  var messageText string
  var messageCounter int
  var messageLimiter int

	// connect to nats server.
	natsConnection, connectionError := nats.Connect(nats.DefaultURL)

  // handle connection error.
  if connectionError != nil {
    // print error message.
    fmt.Printf("Connection to NATS server '%s' failed. '%s'.\n", nats.DefaultURL, connectionError.Error())

    // exit with error.
    os.Exit(1)
  }

  // add task to close the connection in the end.
  defer natsConnection.Close()

  // print message.
  fmt.Printf("Connected to '%s'.\n", nats.DefaultURL)

	// set nats topic name.
	subject = "foo"

  // set message text.
  messageText = "Hello World!"

  // initialize message counter and limit.
  messageCounter = 0;
  messageLimiter = 100;

  // record start time.
  start := time.Now()

  for messageCounter < messageLimiter {
    // increment message counter.
    messageCounter++

    // set message text.
    messageText = "Something cool" + strconv.Itoa(messageCounter)

    // publish message.
    natsConnection.Publish(subject, []byte(messageText))

    // print message; disabled to avoid unnecessary overhead to the measurements.
    //fmt.Printf("Published message '%s' to subject '%s'.\n", messageText, subject)
  }

  // record stop time.
  stop := time.Now()

  // calculate elapsed time.
  elapsed := stop.Sub(start)

  // print results.
  fmt.Println("Elapsed time:")
  fmt.Println(elapsed)
  fmt.Println("Time used per transaction:")
  fmt.Println(elapsed.Seconds() / float64(messageLimiter))
  fmt.Println("Transactions per second:")
  fmt.Println(float64(messageLimiter) / elapsed.Seconds())

}
