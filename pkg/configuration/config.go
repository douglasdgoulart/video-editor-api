package configuration

type KafkaProducerConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}
