package main

import (
  "log"

  "github.com/Yousiph1/imageResizer/consumer/utils"
  "github.com/streadway/amqp"
)



func main() {

conn, err := amqp.Dial("amqp://localhost:5672")
utils.FailOnError(err, "Failed to connect to rabbitmq server")

ch, err := conn.Channel()
        utils.FailOnError(err, "Failed to open a channel")
        defer ch.Close()

q, err := ch.QueueDeclare(
               "image-queue", // name
               true,         // durable
               false,        // delete when unused
               false,        // exclusive
               false,        // no-wait
               nil,          // arguments
       )
     utils.FailOnError(err, "Failed to declare a queue")

images, err := ch.Consume(
              q.Name, // queue
              "",     // consumer
              false,  // auto-ack
              false,  // exclusive
              false,  // no-local
              false,  // no-wait
              nil,    // args
      )
     utils.FailOnError(err, "Failed to register a consumer")


forever := make(chan bool)

 go func ()  {
 for image := range images {
    fileName := string(image.Body)
    go utils.Resize(fileName, "large")
    go utils.Resize(fileName, "medium")
    go utils.Resize(fileName, "small")
    image.Ack(false)
  }
 }()

log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
<- forever
}
