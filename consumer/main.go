package main

import "github.com/Shopify/sarama"

func main() {

	server := []string{"localhost:9092"}

	consumer, err := sarama.NewConsumer(server, nil)

	if err != nil {
		panic(err)
	}

	consumer.Close()
	defer 


}
