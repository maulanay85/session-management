package config

import "github.com/nsqio/go-nsq"

type NSQClient struct {
	Producer *nsq.Producer
	Consumer *nsq.Consumer
}

func InitializeNSQ(conf Config) (NSQClient, error) {
	producer, err := nsq.NewProducer(conf.NsqUrl, nsq.NewConfig())
	if err != nil {
		return NSQClient{}, err
	}

	// Initialize consumer (no subscription yet)
	// consumer, err := nsq.NewConsumer("audit_trail", "audit_trail", nsq.NewConfig())
	// if err != nil {
	// 	return NSQClient{}, err
	// }

	return NSQClient{
		Producer: producer,
	}, nil

}
