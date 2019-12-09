package kafka_test

import (
	"context"

	"github.com/Shopify/sarama"
)

type adminClient struct {
	ListTopicsFunction func() (map[string]sarama.TopicDetail, error)
}

func (c *adminClient) ListTopics() (map[string]sarama.TopicDetail, error) {
	return c.ListTopicsFunction()
}

type consumer struct {
	PartitionsFunction       func(topic string) ([]int32, error)
	ConsumePartitionFunction func(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error)
}

func (c *consumer) ConsumePartition(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
	return c.ConsumePartitionFunction(topic, partition, offset)
}

func (c *consumer) Partitions(topic string) ([]int32, error) {
	return c.PartitionsFunction(topic)
}

type partitionConsumer struct {
	messages  chan *sarama.ConsumerMessage
	wasClosed bool
}

func (c *partitionConsumer) AsyncClose() {

}

func (c *partitionConsumer) Close() error {
	c.wasClosed = true
	return nil
}

func (c *partitionConsumer) Errors() <-chan *sarama.ConsumerError {
	return nil
}

func (c *partitionConsumer) HighWaterMarkOffset() int64 {
	return 0
}

func (c *partitionConsumer) Messages() <-chan *sarama.ConsumerMessage {
	return c.messages
}

type decoder struct {
	DecodeFunction func(ctx context.Context, id uint32, message []byte) (string, error)
}

func (d *decoder) Decode(ctx context.Context, id uint32, message []byte) (string, error) {
	return d.DecodeFunction(ctx, id, message)
}
