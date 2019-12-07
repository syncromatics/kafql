package kafka_test

import (
	"context"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/syncromatics/kafql/internal/kafka"
	"github.com/syncromatics/kafql/internal/models"

	"github.com/stretchr/testify/assert"
)

func Test_Service_GetTopics(t *testing.T) {
	adminClient := &adminClient{}
	consumerClient := &consumer{}
	decoder := &decoder{}

	adminClient.ListTopicsFunction = func() (map[string]sarama.TopicDetail, error) {
		return map[string]sarama.TopicDetail{
			"stuff": sarama.TopicDetail{},
		}, nil
	}

	service := kafka.NewService(adminClient, consumerClient, decoder)

	topics, err := service.GetTopics()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []*models.Topic{
		&models.Topic{
			Name: "stuff",
		},
	}, topics)
}

func Test_Service_SubscribePlainText(t *testing.T) {
	adminClient := &adminClient{}
	consumerClient := &consumer{}
	decoder := &decoder{}

	consumerClient.PartitionsFunction = func(topic string) ([]int32, error) {
		return []int32{0, 1}, nil
	}

	service := kafka.NewService(adminClient, consumerClient, decoder)

	ctx, cancel := context.WithCancel(context.Background())

	fromKafka := make(chan *sarama.ConsumerMessage, 1)
	consumerClient.ConsumePartitionFunction = func(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
		assert.Equal(t, "Test_Service_Subscribe", topic)
		assert.Equal(t, int32(0), partition)
		assert.Equal(t, int64(-1), offset)
		return &partitionConsumer{messages: fromKafka}, nil
	}

	r, err := service.Subscribe(ctx, "Test_Service_Subscribe", "Test_Service_Subscribe_1")
	if err != nil {
		t.Fatal(err)
	}

	fromKafka <- &sarama.ConsumerMessage{
		Key:   []byte("Test_Service_Subscribe_1"),
		Value: []byte("just plain text"),
	}

	timer := time.NewTimer(10 * time.Second)
	select {
	case m := <-r:
		assert.Equal(t, &models.Record{
			Key:   "Test_Service_Subscribe_1",
			Value: "just plain text",
		}, m)

	case <-timer.C:
		t.Fatal("timeout")
	}

	cancel()
}

func Test_Service_SubscribeProto(t *testing.T) {
	adminClient := &adminClient{}
	consumerClient := &consumer{}
	decoder := &decoder{}

	consumerClient.PartitionsFunction = func(topic string) ([]int32, error) {
		return []int32{0, 1}, nil
	}

	service := kafka.NewService(adminClient, consumerClient, decoder)

	ctx, cancel := context.WithCancel(context.Background())

	fromKafka := make(chan *sarama.ConsumerMessage, 1)
	consumerClient.ConsumePartitionFunction = func(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
		assert.Equal(t, "Test_Service_SubscribeProto", topic)
		assert.Equal(t, int32(0), partition)
		assert.Equal(t, int64(-1), offset)
		return &partitionConsumer{messages: fromKafka}, nil
	}

	decoder.DecodeFunction = func(ctx context.Context, id uint32, message []byte) (string, error) {
		assert.Equal(t, uint32(9), id)
		assert.Equal(t, []byte{0x45, 0x67}, message)
		return "decoded", nil
	}

	r, err := service.Subscribe(ctx, "Test_Service_SubscribeProto", "Test_Service_SubscribeProto_2")
	if err != nil {
		t.Fatal(err)
	}

	fromKafka <- &sarama.ConsumerMessage{
		Key:   []byte("Test_Service_SubscribeProto_2"),
		Value: []byte{0x2, 0x0, 0x0, 0x0, 0x9, 0x45, 0x67},
	}

	timer := time.NewTimer(10 * time.Second)
	select {
	case m := <-r:
		assert.Equal(t, &models.Record{
			Key:   "Test_Service_SubscribeProto_2",
			Value: "decoded",
		}, m)

	case <-timer.C:
		t.Fatal("timeout")
	}

	cancel()
}
