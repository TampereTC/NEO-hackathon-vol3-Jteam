package main

/*
  The original application is in https://nats.io/documentation/tutorials/nats-client-dev/,
  which has been modified and extended for use of struct/JSON, byte encoder, parallel
  execution, documentation and performace measurements.
*/

// import packages.
import "fmt"
import "log"
import "os"
import "time"
import "sync"
import "encoding/json"
import "github.com/nats-io/go-nats"
import "./message"

// main function.
func main() {
  // variables.
  var publishTopic string
  var numberOfMessages int
  var numberOfThreads int
  var threadIndex int
  var threads sync.WaitGroup
  var messagesPerSecond float64

  // set number of messages, number of threads and thread index.
  numberOfMessages = 1000000
  numberOfThreads = 200
  threadIndex = 0

  // initialize threads.
  threads.Add(numberOfThreads)

  // set nats publish topic name.
	publishTopic = "foo"

  // record start time.
  start := time.Now()

  // go through all threads.
  for threadIndex < numberOfThreads {
    // execute a thread.
    go publishMessage(publishTopic, threadIndex, numberOfMessages/numberOfThreads, &threads)

    // print log message.
    //log.Printf("Thread %d executed.\n", threadIndex)

    // increment thread index.
    threadIndex++
  }

  // wait until all threads are completed.
  threads.Wait()

  // record stop time.
  stop := time.Now()

  // calculate elapsed time.
  elapsed := stop.Sub(start)

  // calculate messages per second.
  messagesPerSecond = float64(numberOfMessages) / elapsed.Seconds()

  // print results.
  fmt.Printf("Number of threads: %d.\n", numberOfThreads)
  fmt.Printf("Number of messages: %d.\n", numberOfMessages)
  fmt.Printf("Elapsed time: %s.\n", elapsed)
  fmt.Printf("Messages per second: %f.\n", messagesPerSecond)
}

// publish message.
func publishMessage(publishTopic string, threadIndex int, numberOfMessages int, threads *sync.WaitGroup) {
  // variables.
  var sentMessage message.Message
  var jsonMessage []byte
  var messageIndex int

	// connect to nats server.
	natsConnection, connectionError := nats.Connect(nats.DefaultURL)

  // handle connection error.
  handleConnectionError(connectionError)

  // add task to close the connection in the end.
  defer natsConnection.Close()

  // print message.
  //log.Printf("Connected to '%s'.\n", nats.DefaultURL)

	// initialize message.
  sentMessage.MessageText = "Lorem ipsum dolor sit amet, consectetur adipisci elit, sed eiusmod tempor incidunt ut labore et dolore magna aliqua."

  // set message index.
  messageIndex = 0;

  // print log message.
  //log.Printf("Publishing messages to '%s'.\n", publishTopic)

  // go through all messages.
  for messageIndex < numberOfMessages {
    // set message counter.
    sentMessage.MessageCounter = messageIndex

    // set last message status.
    if messageIndex == numberOfMessages {
      sentMessage.LastMessage = true
    } else {
      sentMessage.LastMessage = false
    }

    // marshall message.
    jsonMessage, _ = json.Marshal(sentMessage)

    // publish message.
    natsConnection.Publish(publishTopic, jsonMessage)

    // increment message counter.
    messageIndex++
  }

  // set thread done.
  threads.Done()
}

// handle connection error.
func handleConnectionError(connectionError error) {
  // error has occured.
  if connectionError != nil {
    // print error message.
    log.Printf("Connection to NATS server '%s' failed. '%s'.\n", nats.DefaultURL, connectionError.Error())

    // exit with error.
    os.Exit(1)
  }
}
