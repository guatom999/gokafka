package services

import "github.com/Shopify/sarama"

type ConsumerHandler struct {
	eventHandler EventHandler
}

func NewConsumerHandler(eventHandler EventHandler) sarama.ConsumerGroupHandler {
	return ConsumerHandler{eventHandler}
}

func (h ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil

}

func (h ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil

}

func (h ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.eventHandler.Handle(msg.Topic, msg.Value)
		session.MarkMessage(msg, "")
	}

	return nil
}
