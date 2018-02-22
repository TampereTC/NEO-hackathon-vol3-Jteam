package main

/*
  The original application is in https://nats.io/documentation/tutorials/nats-client-dev/,
  and it has been used almost as is. I've extended some variable names and added comments
  to clarify my own understanding.

  Ping-Pong edition: producer sends all messages, and waits a notification from consumer,
  when all messages have been processed.
*/

// import packages.
import "log"
import "os"
import "github.com/nats-io/go-nats"
import "time"
import "encoding/json"
import "./message"

// main function.
func main() {
  // variables.
  var publishTopic string
  //var subscribeTopic string
  var sentMessage message.Message
  var messageCounter int
  var messageLimiter int
  var jsonMessage []byte
  //var messageAlley = make(chan []byte, 1)

	// connect to nats server.
	natsConnection, connectionError := nats.Connect(nats.DefaultURL)

  // handle connection error.
  if connectionError != nil {
    // print error message.
    log.Printf("Connection to NATS server '%s' failed. '%s'.\n", nats.DefaultURL, connectionError.Error())

    // exit with error.
    os.Exit(1)
  }

  // add task to close the connection in the end.
  defer natsConnection.Close()

  // print message.
  log.Printf("Connected to '%s'.\n", nats.DefaultURL)

	// set nats publish topic name.
	publishTopic = "foo"

  // initialize message.
  sentMessage.MessageText = "Lorem ipsum dolor sit amet, consectetur adipisci elit, sed eiusmod tempor incidunt ut labore et dolore magna aliqua."

  // initialize message counter and limit.
  messageCounter = 0;
  messageLimiter = 1000000;

  // print log message.
  log.Printf("Publishing to '%s'.\n", publishTopic)

  // record start time.
  start := time.Now()

  for messageCounter < messageLimiter {
    // increment message counter.
    messageCounter++

    // set message counter.
    sentMessage.MessageCounter = messageCounter

    // set last message status.
    if messageCounter == messageLimiter {
      sentMessage.LastMessage = true
    } else {
      sentMessage.LastMessage = false
    }

    // marshall message.
    jsonMessage, _ = json.Marshal(sentMessage)

    // publish message.
    natsConnection.Publish(publishTopic, jsonMessage)
  }

/*
  // set nats subscribe topic name.
  //subscribeTopic = "bar"

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
          fmt.Println("Message received")
          // unmarshal json; []byte.
          json.Unmarshal(receivedJson, &receivedMessage)

          // last message received.
          if receivedMessage.LastMessage == true {
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
        // default action.
        default:
          // do nothing.
      }
    }
  }()

  // print log message.
  fmt.Printf("Subscribing from '%s'.\n", subscribeTopic)

  // subscribe messages.
	natsConnection.Subscribe(subscribeTopic, func(message *nats.Msg) {
		// forward message data for processing.
    fmt.Println(message)
    //messageAlley <- message.Data
	})
*/
  // record stop time.
  stop := time.Now()

  // calculate elapsed time.
  elapsed := stop.Sub(start)

  // print results.
  log.Println("Elapsed time:")
  log.Println(elapsed)
  log.Println("Time used per transaction:")
  log.Println(elapsed.Seconds() / float64(messageLimiter))
  log.Println("Transactions per second: ")
  log.Println(float64(messageLimiter) / elapsed.Seconds())
}
