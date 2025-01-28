package server

import (
	"chat/models"
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
)

var poolSizeClickhouse = 1000

func (c *ClickHouse) StartConsumerDebeziumChats() {
	ctx := context.Background()

	messageChan := make(chan *kafka.Message, poolSizeClickhouse)
	workerPool := make(chan struct{}, poolSizeClickhouse)

	go func() {
		defer c.clickHouseCli.Client.Close()

		for {
			msg, err := c.clickHouseCli.KafkaConsumerDebeziumChats.FetchMessage(ctx)

			if err != nil {
				log.Println("[ERROR] Потребитель debezium: %s", err.Error())

				if ctx.Err() != nil {
					close(messageChan)
					log.Println("[ERROR] Потребитель debezium ctx: %s", err.Error())
					return
				}

				continue
			}

			messageChan <- &msg
		}
	}()

	for i := 0; i < poolSizeClickhouse; i++ {
		go func() {
			for msg := range messageChan {
				workerPool <- struct{}{}
				c.addMessageDebeziumChats(msg, ctx)
				<-workerPool
			}
		}()
	}
}

func (c *ClickHouse) addMessageDebeziumChats(msg *kafka.Message, ctx context.Context) {
	if msg == nil {
		return
	}

	chat := new(models.ChatsChanges)
	afterAndBefore := new(models.AfterAndBefore)

	err := jsonUnmarshalAfterAndBefore(msg.Value, afterAndBefore, chat)
	if err != nil {
		log.Println("[ERROR]", err)
		return
	}

	err = c.debeziumChat.AddMessageDebeziumChats(chat)
	if err != nil {
		log.Println("[ERROR] AddMessage потребитель debezium", err)
		return
	}

	if err = c.clickHouseCli.KafkaConsumerDebeziumChats.CommitMessages(ctx, *msg); err != nil {
		log.Println("Не удалось зафиксировать сообщение debezium:", err)
	}
}

func jsonUnmarshalAfterAndBefore(msg []byte, afterAndBefore *models.AfterAndBefore, chat *models.ChatsChanges) error {
	err := json.Unmarshal(msg, &afterAndBefore.After)
	if err != nil {
		return errors.New(" Unmarshal потребитель debezium after " + err.Error())
	}

	if afterAndBefore.After.After != "" {
		err = json.Unmarshal([]byte(afterAndBefore.After.After), &chat.After)
		if err != nil {
			return errors.New("Unmarshal потребитель debezium chat after " + err.Error())
		}
	}

	err = json.Unmarshal(msg, &afterAndBefore.Before)
	if err != nil {
		return errors.New("Unmarshal потребитель debezium before " + err.Error())
	}

	if afterAndBefore.Before.Before != "" {
		err = json.Unmarshal([]byte(afterAndBefore.Before.Before), &chat.Before)
		if err != nil {
			return errors.New("Unmarshal потребитель debezium chat before " + err.Error())
		}
	}

	err = json.Unmarshal(msg, &afterAndBefore.Op)
	if err != nil {
		return errors.New(" Unmarshal потребитель debezium op " + err.Error())
	}

	chat.Op = afterAndBefore.Op.Op

	return nil
}
