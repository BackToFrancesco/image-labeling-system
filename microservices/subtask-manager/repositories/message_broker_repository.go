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

func (m MessageBrokerRepository) PublishNewSubtask(subtask *models.SubtaskMessage) error {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)

	err := encoder.Encode(subtask)
	if err != nil {
		return err
	}

	err = m.broker.Channel.PublishWithContext(context.Background(),
		"",
		datasources.NewSubtasksQueue,
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

func (m MessageBrokerRepository) ConsumeCompletedSubtasks(consume func(message *models.CompletedSubtaskMessage) error) {
	messages, err := m.broker.Channel.Consume(datasources.CompletedSubtasksQueue,
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
				var completedSubtaskMessage models.CompletedSubtaskMessage

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
