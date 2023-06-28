package main

import (
	"consumer/repositories"
	"consumer/services"
	"context"
	"events"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func initData() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.database"),
	)

	dial := mysql.Open(dsn)

	db, err := gorm.Open(dial, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(err)
	}

	return db

}

func main() {

	consumerGroup, err := sarama.NewConsumerGroup(viper.GetStringSlice("kafka.servers"), viper.GetString("kafka.group"), nil)

	if err != nil {
		panic(err)
	}

	defer consumerGroup.Close()

	db := initData()

	accountRepo := repositories.NewAccountRepositoryDB(db)
	accountEventHandler := services.NewAccountEventHandler(accountRepo)
	accountConsumerHandler := services.NewConsumerHandler(accountEventHandler)

	fmt.Println("Account Consumer Start")

	for {
		consumerGroup.Consume(context.Background(), events.Topics, accountConsumerHandler)
	}

	// consumerGroup.Consume()

	// var db *gorm.DB

	// chon := repositories.NewAccountRepositoryDB(db)

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
