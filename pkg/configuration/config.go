package configuration

type KafkaProducerConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
}

type KafkaConsumerConfig struct {
	Brokers []string `json:"brokers"`
	GroupID string   `json:"group_id"`
	Topic   string   `json:"topic"`
	Offset  string   `json:"offset"`
}
