package main

/*
The original application is in https://github.com/vsouza, which has been modified
and extended for use of struct/JSON, byte encoder, parallel execution, documentation
and performace measurements.
*/

// import packages.
import "fmt"
//import "log"
import "os"
import "time"
import "sync"
import "encoding/json"
import "github.com/Shopify/sarama"
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
  numberOfThreads = 5
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
    //log.Printf("Thread %d executes.\n", threadIndex)

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
  var brokers []string
  var sentMessage message.Message
  var jsonMessage []byte
  var messageIndex int

  // initialize configuration.
	config := sarama.NewConfig()

  // set configuration.
  config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

  // set brokers.
  brokers = []string{"localhost:9092"}

	// initialize producer.
  producer, producerInitializationError := sarama.NewSyncProducer(brokers, config)

  // handle producer initialization error.
  if producerInitializationError != nil {
		// print error message.
    fmt.Printf("Producer initialization failed. '%s'.\n", producerInitializationError.Error())

    // exit with error.
    os.Exit(1)
	}

  // add taks to close the producer in the end.
  defer func() {
    // close producer.
    producerClosingError := producer.Close()

    // handle producer closing error.
    if producerClosingError != nil {
      // print error message.
      fmt.Printf("Producer closing failed. '%s'.\n", producerClosingError.Error())

      // exit with error.
      os.Exit(1)
		}
	}()

  // print message.
  //log.Printf("Connected to 'localhost:9092'.\n")

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

    // create message.
    message := &sarama.ProducerMessage{
      Topic: publishTopic,
      Value: sarama.ByteEncoder(jsonMessage),
    }

    // send message.
    _, _, messageSendingError := producer.SendMessage(message)

    // handle message sending error.
    if messageSendingError != nil {
      // print error message.
      fmt.Printf("Message sending failed. '%s'.\n", messageSendingError.Error())

      // exit with error.
      os.Exit(1)
  	}

    // increment message counter.
    messageIndex++
  }

  // set thread done.
  threads.Done()
}
