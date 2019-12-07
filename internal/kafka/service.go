package kafka

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/syncromatics/kafql/internal/models"

	"github.com/Shopify/sarama"
	"github.com/burdiyan/kafkautil"
	"github.com/pkg/errors"
)

// ClusterAdmin is the kafka cluster admin client
type ClusterAdmin interface {
	ListTopics() (map[string]sarama.TopicDetail, error)
}

// Consumer is the kafka consumer client
type Consumer interface {
	ConsumePartition(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error)
	Partitions(topic string) ([]int32, error)
}

// Service is the kafka service
type Service struct {
	adminClient  ClusterAdmin
	consumer     Consumer
	protoDecoder Decoder
}

// Decoder can decode a message based on the schema id
type Decoder interface {
	Decode(ctx context.Context, id uint32, message []byte) (string, error)
}

// NewService creates a new kafka service
func NewService(admin ClusterAdmin, consumer Consumer, protoDecoder Decoder) *Service {
	return &Service{
		adminClient:  admin,
		consumer:     consumer,
		protoDecoder: protoDecoder,
	}
}

// GetTopics returns all the topics in kafka
func (s *Service) GetTopics() ([]*models.Topic, error) {
	topics, err := s.adminClient.ListTopics()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list topics")
	}

	names := []*models.Topic{}
	for t := range topics {
		topic := &models.Topic{Name: t}
		names = append(names, topic)
	}

	return names, nil
}

// Subscribe subscribes to a kafka topic by key
func (s *Service) Subscribe(ctx context.Context, topic string, key string) (<-chan *models.Record, error) {
	partitions, err := s.consumer.Partitions(topic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get partitions")
	}

	p := kafkautil.NewJVMCompatiblePartitioner(topic)
	partition, err := p.Partition(&sarama.ProducerMessage{Key: sarama.StringEncoder(key)}, int32(len(partitions)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get partition")
	}

	update := make(chan *models.Record)
	consume, err := s.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)

	go func() {
		defer close(update)
		defer consume.Close()
		for {
			select {
			case msg := <-consume.Messages():
				if string(msg.Key) != key {
					continue
				}
				if err != nil {
					fmt.Printf("failed reading message: %v", err)
					return
				}

				if len(msg.Value) == 0 {
					continue
				}

				response := ""

				switch msg.Value[0] {
				case 0x0: // avro

				case 0x2: //proto
					id := int32(binary.BigEndian.Uint32(msg.Value[1:5]))
					response, err = s.protoDecoder.Decode(ctx, uint32(id), msg.Value[5:])
					if err != nil {
						fmt.Printf("failed decoding message: %v", err)
						return
					}

				default: // just string
					response = string(msg.Value)
				}

				select {
				case update <- &models.Record{
					Key:   string(msg.Key),
					Value: response,
				}:

				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return update, nil
}
