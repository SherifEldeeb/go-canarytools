package canarytools

import (
	"crypto/tls"

	"github.com/Shopify/sarama"

	log "github.com/sirupsen/logrus"
)

// KafkaForwarder sends alerts to kafka
type KafkaForwarder struct {
	// w *kafka.Writer
	p     sarama.SyncProducer
	topic string
	l     *log.Logger
	// TODO: TLS!
}

// NewKafkaForwarder creates a new KafkaForwarder
func NewKafkaForwarder(brokers []string, topic string, tlsconfig *tls.Config, l *log.Logger) (kafkaforwarder *KafkaForwarder, err error) {
	kafkaforwarder = &KafkaForwarder{}
	kafkaforwarder.l = l
	kafkaforwarder.topic = topic
	sarama.Logger = l
	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true
	config.ClientID = "CanaryChirpForwarder"
	if tlsconfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsconfig
	}
	kafkaforwarder.p, err = sarama.NewSyncProducer(brokers, config)
	return
}

func (kf KafkaForwarder) Forward(outChan <-chan []byte, incidentAckerChan chan<- []byte) {
	for i := range outChan {
		kf.l.WithFields(log.Fields{
			"source": "KafkaForwarder",
			"stage":  "forward",
		}).Debug("Kafka out incident")
		_, _, err := kf.p.SendMessage(&sarama.ProducerMessage{
			Topic: kf.topic,
			Value: sarama.ByteEncoder(i),
		})
		if err != nil {
			kf.l.WithFields(log.Fields{
				"source": "KafkaForwarder",
				"stage":  "forward",
				"err":    err,
			}).Error("Kafka error")
		}

		incidentAckerChan <- i
	}
}
