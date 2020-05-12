package kafkautil

import (
	"github.com/Shopify/sarama"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	consumer, err := NewConsumer(addrs, "test-group-id", []string{"test-topic"}, sarama.OffsetNewest, true)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		msg, err := consumer.Receive(NeverTimeout)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("topic:%v, partition:%v/%v, %v", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		time.Sleep(2 * time.Second)
	}
}
