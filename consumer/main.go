package main

import (
	"consumer/repositories"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	consumerGroup, err := sarama.NewConsumerGroup(viper.GetStringSlice("kafka.server"), viper.GetString("kafka.gro"), nil)

	if err != nil {
		panic(err)
	}

	var db *gorm.DB

	chon := repositories.NewAccountRepositoryDB(db)

	defer consumerGroup.Close()

	// consumerGroup.Consume()

}

// func main() {

// 	server := []string{"localhost:9092"}

// 	consumer, err := sarama.NewConsumer(server, nil)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer consumer.Close()

// 	partitionConsumer, err := consumer.ConsumePartition("bosshello", 0, sarama.OffsetNewest)

// 	if err != nil {

// 	defer partitionConsumer.Close()

// 	fmt.Print("Consumer start.")

// 	for {
// 		select {
// 		case err := <-partitionConsumer.Errors():
// 			fmt.Println(err)
// 		case msg := <-partitionConsumer.Messages():
// 			fmt.Println(string(msg.Value))
// 		}
// 	}

// }
