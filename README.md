# NEO-hackathon-vol3-Jteam

#; i.e. ping-pong edition.# Kafka and Go

Kafka message producer and consumer implemented with Go:
- Kafka-Serial-Producer.go
- Kafka-Serial-Consumer.go

The base code for the applications used in this experimentation are copied from existing Github projects. Primarily the source code is used as is. Some of the variable names are extended and new comments are added to improve my own understanding.

The source of the original code is documented in the individual applications.

Time recording points are added to the code to enable measurements.

## NATS and Go

NATS message producer and consumer implemented with Go:
* NATS-Serial-Producer.go
* NATS-Serial-Consumer.go
* NATS-Serial-PP-Producer.go; i.e. ping-pong edition.
* NATS-Serial-PP-Consumer.go; i.e. ping-pong edition.

The base code for the applications used in this experimentation are copied from NATS web site. Primarily the source code is used as is. Some of the variable names are extended and new comments are added to improve my own understanding.

Time recording points are added to the code to enable measurements.

The source of the original code is documented in the individual applications.

In the ping-pong edition the consumer will send message back, when all messages have been consumed. The produced message is a JSON ([]byte), which contains a boolean value for last message ("last_message": true).
