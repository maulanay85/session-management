package consumer

import (
	"context"
	"encoding/json"
	"log"
	"scs-session/internal/domain"
	"scs-session/internal/usecase"

	"github.com/nsqio/go-nsq"
)

type NSQConsumer struct {
	auditTrailUsecase usecase.AuditTrailUsecase
}

func (c *NSQConsumer) HandleMessageAuditTrail(m *nsq.Message) error {
	var auditTrail domain.AuditTrail

	err := json.Unmarshal(m.Body, &auditTrail)
	if err != nil {
		log.Printf("error unmarshal message %v", err)
		return err
	}

	err = c.auditTrailUsecase.HandleAuditTrailMessage(context.Background(), auditTrail)
	if err != nil {
		log.Printf("error on HandleAuditTrailMessage: %+v", err)
		return err
	}
	return nil
}

func NewNSQConsumer(at usecase.AuditTrailUsecase) *NSQConsumer {
	return &NSQConsumer{
		auditTrailUsecase: at,
	}
}

func StartNSQConsumer(nsqAddress string, topic string, consumer *NSQConsumer) error {
	cfg := nsq.NewConfig()
	nsqConsumer, err := nsq.NewConsumer(topic, "channel", cfg)

	if err != nil {
		return err
	}

	nsqConsumer.AddHandler(nsq.HandlerFunc(consumer.HandleMessageAuditTrail))

	err = nsqConsumer.ConnectToNSQD(nsqAddress)
	if err != nil {
		return err
	}

	// Block until the consumer stops
	<-nsqConsumer.StopChan
	return nil
}
