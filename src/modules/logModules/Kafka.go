package logModules

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"goshort/src/kernel/utils"
	"goshort/src/kernel/utils/other"
	"goshort/src/types"
	errors2 "goshort/src/types/errors"
	"io"
	"net"
	"strconv"
	"sync"
)

type Kafka struct {
	types.LoggerBase
	conn      *kafka.Conn
	ip        string
	port      int
	topic     string
	partition int
	Kernel    types.KernelInterface
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

func CreateKafka(kernel types.KernelInterface) types.LoggerInterface {
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

func (k *Kafka) Run(wg *sync.WaitGroup) error {
	wg.Done()
	err := k.Connect()
	if err != nil {
		return err
	}
	k.LoggerBase.ModuleBase.IsAvailableVal = 1
	k.Kernel.SetModuleRunState(k)
	return nil
}

func (k *Kafka) Stop() error {
	k.SetUnavailableAndTryGetReconnectionControl()
	k.SetDeath()
	if k.conn != nil {
		err := k.conn.Close()
		k.Kernel.SetModuleStopState(k)
		return err
	}
	return nil
}

func (k *Kafka) Send(element types.Log) error {
	if !k.IsAvailable() {
		return nil
	}
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

func (k *Kafka) SendError(err error) error {
	return k.Send(other.InterfaceToLogWrapper(err))
}

func (k *Kafka) SendBatch(batch *types.LoggingQueueNode) error {
	if !k.IsAvailable() {
		return nil
	}

	var messages []kafka.Message

	for batch != nil {
		data, _ := json.Marshal(batch.Log.ToMap())
		messages = append(messages, kafka.Message{Value: data})
		batch = batch.Next
	}

	_, err := k.conn.WriteMessages(messages...)
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

	return nil
}

func (k *Kafka) GetName() string {
	return k.Name
}

func (k *Kafka) GetType() string {
	return "Logger.Kafka"
}

func (k *Kafka) TryReconnect() error {
	return k.Connect()
}
