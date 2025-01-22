package repository

import (
	"encoding/json"
	"log"
	"scs-session/internal/config"
	"scs-session/internal/domain"

	"github.com/nsqio/go-nsq"
)

type NSQRepositoryImpl struct {
	nsq *config.NSQClient
}

// PublishMessage implements NSQRepository.
func (r NSQRepositoryImpl) PublishMessage(data domain.AuditTrail, topic string) error {
	messageData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Publish to NSQ
	err = r.nsq.Producer.Publish(topic, messageData)
	if err != nil {
		return err
	}
	log.Printf("Message sent to topic: %s", topic)
	return nil
}

type NSQRepository interface {
	PublishMessage(data domain.AuditTrail, topic string) error
}

func (r *NSQRepositoryImpl) HandleAuditTrailMessage(m *nsq.Message) error {
	// var auditTrail domain.AuditTrail
	// err := json.Unmarshal(m.Body, &auditTrail)
	// r.nsq.Consumer.han
	return nil
}

func NewNSQRepository(nsq *config.NSQClient) NSQRepository {
	return NSQRepositoryImpl{
		nsq: nsq,
	}
}
