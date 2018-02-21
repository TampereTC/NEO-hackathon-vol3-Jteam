package main

/*
  The original application is in https://github.com/Shopify/sarama, and it has
  been used almost as is. I've extended some variable names and added comments
  to clarify my own understanding.
*/

// import packages.
import "fmt"
import "os"
import "os/signal"
import "github.com/Shopify/sarama"

// main function.
func main() {
  // variables.
  var brokers []string
  var topic string

  // initialize configuration.
  config := sarama.NewConfig()

  // set configuration.
	config.Consumer.Return.Errors = true

	// set brokers.
	brokers = []string{"localhost:9092"}

	// initialize consumer.
	consumer, consumerInitializationError := sarama.NewConsumer(brokers, config)

  // handle consumer initialization error.
  if consumerInitializationError != nil {
    // print error message.
    fmt.Printf("Consumer initialization failed. '%s'.\n", consumerInitializationError.Error())

    // exit with error.
    os.Exit(1)
	}

  // add taks to close the consumer in the end.
	defer func() {
    // close consumer.
    consumerClosingError := consumer.Close()

    // handle consumer closing error.
    if consumerClosingError != nil {
      // print error message.
      fmt.Printf("Consumer closing failed. '%s'.\n", consumerClosingError.Error())

      // exit with error.
      os.Exit(1)
    }
  }()

  // set kafka topic name.
	topic = "test"

  // create partition consumer for "test" topic. The messages are consumed
  // from partition zero, and at start up all the existing messages shall
  // be consumed.
  //
  // TODO figure out how the the consumer can know the number of partitions?
	partitionConsumer, consumerCreationError := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)

  // handle consumer creation error.
  if consumerCreationError != nil {
    // print error message.
    fmt.Printf("Consumer creation failed. '%s'.\n", consumerCreationError.Error())

    // exit with error.
    os.Exit(1)
	}

  // create channel for operating system signals.
  signals := make(chan os.Signal, 1)

  // add interruption (os.Interrupt) as part of monitored signals.
  signal.Notify(signals, os.Interrupt)

	// initialize message counter.
	messageCounter := 0

	// create non-buffered channel type of emtpty struct.
	doneCh := make(chan struct{})

  // execute go routine.
  go func() {
    // print message.
    fmt.Println("Start listening")

    // start end-less loop.
    for {
      // wait for input.
      select {
        // error is received.
        case err := <-partitionConsumer.Errors():
          // print error message.
          fmt.Printf("Received error. '%s'.\n", err.Error())
        // message is received.
        case msg := <-partitionConsumer.Messages():
          // increase message count.
          messageCounter++

          // print message.
          fmt.Printf("Received message. %s %s.\n", string(msg.Key), string(msg.Value))
        // interruption is received.
        case <-signals:
          // print message.
          fmt.Println("Received interruption.")

          // TODO figure out what is happening here and why.
          doneCh <- struct{}{}
      }
    }
  }()

  // TODO figure out what is happening here and why.
  <-doneCh

  // print message.
  fmt.Println("Processed", messageCounter, "messages")
}
