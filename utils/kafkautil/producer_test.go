package kafkautil

import (
	"github.com/Shopify/sarama"
	"testing"
	"time"
)

var addrs = []string{"148.70.248.33:9092"}

func TestAsyncProducer(t *testing.T) {
	producer, err := NewAsyncProducer(addrs, sarama.WaitForAll, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		s := time.Now().String()
		err := producer.Send(&sarama.ProducerMessage{
			Topic: "test-topic",
			Key:   sarama.StringEncoder(s),
			Value: sarama.StringEncoder(s),
		}, 3*time.Second)
		if err != nil {
			t.Fatal(err)
		}

		v := map[string]string{
			"json-value": s,
		}
		err = producer.SendJSON("test-topic", s, &v, NeverTimeout)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("produce msg, i:%d", i)
		time.Sleep(3 * time.Second)
	}

}

func TestSyncProducer(t *testing.T) {
	producer, err := NewSyncProducer(addrs, sarama.WaitForLocal)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		s := time.Now().String()
		_, _, err := producer.Send(&sarama.ProducerMessage{
			Topic: "test-sync-topic",
			Key:   sarama.StringEncoder(s),
			Value: sarama.StringEncoder(s),
		})
		if err != nil {
			t.Fatal(err)
		}

		v := map[string]string{
			"json-value": s,
		}
		_, _, err = producer.SendJSON("test-topic", s, &v)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("produce msg, i:%d", i)
		time.Sleep(3 * time.Second)
	}
}
