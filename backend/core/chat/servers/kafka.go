package servers

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

var poolSize = 1000

func (h *HttpServer) startConsumer() {
	ctx := context.Background()

	messageChan := make(chan *kafka.Message, poolSize)
	workerPool := make(chan struct{}, poolSize)

	go func() {
		defer h.kafkaConsumer.Close()

		for {
			msg, err := h.kafkaConsumer.FetchMessage(ctx)
			log.Println("FetchMessage", msg)
			if err != nil {
				log.Println("[ERROR] Потребитель: %s", err.Error())

				if ctx.Err() != nil {
					close(messageChan)
					log.Println("[ERROR] Потребитель ctx: %s", err.Error())
					return
				}

				continue
			}

			messageChan <- &msg
		}
	}()

	for i := 0; i < poolSize; i++ {
		go func() {
			for msg := range messageChan {
				workerPool <- struct{}{}
				h.addMessage(msg, ctx)
				<-workerPool
			}
		}()
	}
}

func (h *HttpServer) addMessage(msg *kafka.Message, ctx context.Context) {
	if msg == nil {
		return
	}

	messageStruct := new(Kafka)

	err := json.Unmarshal(msg.Value, messageStruct)
	if err != nil {
		log.Println("[ERROR] Unmarshal потребитель", err)
	}

	messageResponse, file, _, UserTo, err := h.manChat.AddMessage(messageStruct.UserFrom, messageStruct.Message)
	if err != nil {
		log.Println("[ERROR] AddMessage потребитель", err)
	}

	//if messageResponse == PING {
	//	h.pingPong(Users[messageStruct.UserFrom], statusUser)
	//}

	messageWithUserTo, err := json.Marshal(`{"Message":"` + messageResponse + `","UserTo":"` + messageStruct.UserFrom + `","File":"` + file + `"}`)
	if err != nil {
		log.Println("[ERROR] Marshal потребитель", err)
	}

	if _, ok := Users[UserTo]; ok {
		Users[UserTo].UChannel <- messageWithUserTo
	}

	if err := h.kafkaConsumer.CommitMessages(ctx, *msg); err != nil {
		log.Println("Не удалось зафиксировать сообщение:", err)
	}
}
