package main

/*
  The original application is in https://github.com/vsouza, and it has been used
  almost as is. I've extended some variable names and added comments to clarify
  my own understanding.
*/

// import packages.
import "fmt"
import "os"
import "github.com/Shopify/sarama"
import "time"
import "encoding/json"
import "./message"

// main function.
func main() {
  //variables.
  var brokers []string
  var publishTopic string
  var sentMessage message.Message
  var messageCounter int
  var messageLimiter int
  var jsonMessage []byte

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

  // set kafka publish topic name.
  publishTopic = "foo"

  // initialize message.
  sentMessage.MessageText = "Lorem ipsum dolor sit amet, consectetur adipisci elit, sed eiusmod tempor incidunt ut labore et dolore magna aliqua."

  // initialize message counter and limit.
  messageCounter = 0;
  messageLimiter = 100000;

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

    // create message.
    message := &sarama.ProducerMessage{
      Topic: publishTopic,
      Value: sarama.ByteEncoder(jsonMessage),
    }

    // send message.
    //partition, offset, messageSendingError := producer.SendMessage(message)
    _, _, messageSendingError := producer.SendMessage(message)

    // handle message sending error.
    if messageSendingError != nil {
      // print error message.
      fmt.Printf("Message sending failed. '%s'.\n", messageSendingError.Error())

      // exit with error.
      os.Exit(1)
  	}

    // print message; disabled to avoid unnecessary overhead to the measurements.
  	//fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
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
