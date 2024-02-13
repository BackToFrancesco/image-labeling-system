package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fabc.it/subtask-manager/datasources"
	"fabc.it/subtask-manager/domain"
	"fabc.it/subtask-manager/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type MessageBrokerRepository struct {
	broker *datasources.MessageBroker
}

func NewMessageBrokerRepository(broker *datasources.MessageBroker) domain.MessageBrokerService {
	return &MessageBrokerRepository{broker: broker}
}

func (m MessageBrokerRepository) PublishCompletedSubtask(subtask *models.CompletedSubtaskMessage) error {
	// send a message with a completed subtask
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)

	err := encoder.Encode(subtask)
	if err != nil {
		return err
	}

	err = m.broker.Channel.PublishWithContext(context.Background(),
		"",
		datasources.CompletedSubtasksQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b.Bytes(),
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (m MessageBrokerRepository) ConsumeNewSubtasks(consume func(message *models.SubtaskMessage) error) { //passo subtask message
	// Consume New Subtask => create a new subtask
	messages, err := m.broker.Channel.Consume(datasources.NewSubtasksQueue,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for message := range messages {
			go func(body []byte) {
				var completedSubtaskMessage models.SubtaskMessage

				buf := bytes.NewBuffer(body)
				decoder := json.NewDecoder(buf)

				err := decoder.Decode(&completedSubtaskMessage)
				if err != nil {
					log.Print(err)
				}

				err = consume(&completedSubtaskMessage)
				if err != nil {
					log.Print(err)
				}
			}(message.Body)
		}
	}()

	<-forever
}
