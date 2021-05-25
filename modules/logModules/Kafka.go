package logModules

import (
	"encoding/json"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"goshort/kernel"
)

type Kafka struct {
	producer *kafka.Producer
	ip       string
	port     string
	topic    string
	kernel   *kernel.Kernel
}

func (k *Kafka) Init(ip string, port string, topic string) {
	*k = Kafka{ip: ip, port: port, topic: topic}

}

func (k *Kafka) Run() error {
	var err error
	k.producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": k.ip})
	return err
}

func (k *Kafka) Send(element kernel.Log) error {
	data, _ := json.Marshal(element.ToMap())
	return k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)
}
