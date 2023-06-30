package main

import (
	"fmt"
	"producer/controller"
	"producer/services"

	"github.com/Shopify/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

	app := fiber.New()

	producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)

	if err != nil {
		panic(err)
	}

	defer producer.Close()

	eventProducer := services.NewEventProducer(producer)
	accountService := services.NewAccountService(eventProducer)
	accountController := controller.NewAccountController(accountService)

	app.Post("/openAccount", accountController.OpenAccount)
	app.Post("/depositFund", accountController.DepositFund)
	app.Post("/withdrawFund", accountController.WithdrawFund)
	app.Post("/closeAccount", accountController.CloseAccount)

	app.Listen(":8000")

}

// func main() {

// 	server := []string{"localhost:9092"}

// 	producer, err := sarama.NewSyncProducer(server, nil)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer producer.Close()

// 	msg := sarama.ProducerMessage{
// 		Topic: "bosshello",
// 		Value: sarama.StringEncoder("Hello world"),
// 	}

// 	partition, offset, err := producer.SendMessage(&msg)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("partition= %v offset= %v", partition, offset)

// }
