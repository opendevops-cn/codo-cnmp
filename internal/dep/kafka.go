package dep

import (
	"context"
	"fmt"
	"github.com/opendevops-cn/codo-golang-sdk/kafka"
	"time"

	"codo-cnmp/internal/conf"
	"github.com/IBM/sarama"
)

type IKafka interface {
	Close(ctx context.Context)
	SendMessage(ctx context.Context, bytes []byte) error
}

type EmptyKafka struct {
	err error
}

func (x *EmptyKafka) SendMessage(ctx context.Context, message []byte) error {
	return fmt.Errorf("kafka配置错误: %w", x.err)
}

func (x *EmptyKafka) Close(ctx context.Context) {
}

type ErrKafka struct {
	err error
}

func (x *ErrKafka) SendMessage(ctx context.Context, message []byte) error {
	return fmt.Errorf("kafka 发送消息失败: %w", x.err)
}

func (x *ErrKafka) Close(ctx context.Context) {
}

type Kafka struct {
	producer sarama.AsyncProducer
	topic    string
	cleanup  func()
}

func NewKafka(ctx context.Context, bc *conf.Bootstrap) IKafka {
	if bc.KAFKA == nil {
		return &EmptyKafka{
			err: fmt.Errorf("kafka配置为空"),
		}
	}
	if bc.KAFKA.GetADDR() == "" {
		return &EmptyKafka{
			err: fmt.Errorf("kafka addr配置为空"),
		}
	}
	if bc.KAFKA.GetTOPIC() == "" {
		return &EmptyKafka{
			err: fmt.Errorf("kafka topic配置为空"),
		}
	}
	timeout := time.Duration(bc.KAFKA.DialTimeout) * time.Second
	if timeout == 0 {
		timeout = 2 * time.Second // 默认超时时间为 2 秒
	}
	asyncProducer, cleanup, err := kafka.NewProducer(
		kafka.WithBootstrapServers(bc.KAFKA.GetADDR()),
		kafka.WithDialTimeout(uint32(timeout.Seconds())),
	)
	if err != nil {
		return &ErrKafka{
			err: err,
		}
	}
	return &Kafka{
		producer: asyncProducer,
		cleanup:  cleanup,
		topic:    bc.KAFKA.GetTOPIC(),
	}
}

func (x *Kafka) SendMessage(ctx context.Context, message []byte) error {
	x.producer.Input() <- &sarama.ProducerMessage{
		Topic: x.topic,
		Value: sarama.ByteEncoder(message),
	}
	select {
	case <-x.producer.Successes():
		return nil
	case errMsg := <-x.producer.Errors():
		return fmt.Errorf("kafka producer发送消息失败: topic: %s, error: %v", x.topic, errMsg.Err)
	}
}

func (x *Kafka) Close(ctx context.Context) {
	x.cleanup()
}
