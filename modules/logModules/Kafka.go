package logModules

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"goshort/kernel"
	"goshort/kernel/utils"
	"goshort/types"
	errors2 "goshort/types/errors"
	"io"
	"net"
	"strconv"
)

type Kafka struct {
	types.LoggerBase
	conn      *kafka.Conn
	ip        string
	port      int
	topic     string
	partition int
	Kernel    *kernel.Kernel
	Name      string
}

func (k *Kafka) Init(config map[string]interface{}) error {
	_ = k.LoggerBase.Init(config)
	k.Name = utils.UnwrapFieldOrDefault(config, "name", "KafkaLogger").(string)
	k.ip = utils.UnwrapFieldOrDefault(config, "ip", "localhost").(string)
	k.port = utils.UnwrapFieldOrDefault(config, "port", 9092).(int)
	k.topic = utils.UnwrapFieldOrDefault(config, "topic", "goshort").(string)
	k.partition = utils.UnwrapFieldOrDefault(config, "partition", 0).(int)
	return nil
}

func CreateKafka(kernel *kernel.Kernel) types.LoggerInterface {
	return &Kafka{Kernel: kernel}
}

func (k *Kafka) Connect() error {
	var err error
	k.conn, err = kafka.DialLeader(context.Background(), "tcp", k.ip+":"+strconv.Itoa(k.port), k.topic,
		k.partition)
	if err != nil {
		_, ok := err.(*net.OpError)
		if ok {
			return &errors2.BadConnectionError{
				Host:      k.ip,
				Port:      k.port,
				Protocol:  "TCP",
				Retryable: true,
			}
		}
	}
	return err
}

func (k *Kafka) Run() error {
	defer k.Kernel.OperationDone()
	return k.Connect()
}

func (k *Kafka) Stop() error {
	if k.conn != nil {
		return k.conn.Close()
	}
	return nil
}

func (k *Kafka) Send(element types.Log) error {
	data, _ := json.Marshal(element.ToMap())
	_, err := k.conn.Write(data)
	if err != nil {
		_, ok := err.(*net.OpError)
		if errors.Is(err, io.EOF) || ok {
			return &errors2.BadConnectionError{
				Host:      k.ip,
				Port:      k.port,
				Protocol:  "TCP",
				Retryable: true,
			}
		}
	}

	return err
}

func (k *Kafka) GetName() string {
	return k.Name
}

func (k *Kafka) GetType() string {
	return "KafkaLogger"
}

func (k *Kafka) TryReconnect() error {
	return k.Connect()
}
