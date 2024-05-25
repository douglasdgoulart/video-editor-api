package configuration

import (
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	Kafka KafkaConfig `mapstructure:"kafka"`
}

type KafkaConfig struct {
	Enabled             bool                `mapstructure:"enabled"`
	KafkaProducerConfig KafkaProducerConfig `mapstructure:"producer"`
	KafkaConsumerConfig KafkaConsumerConfig `mapstructure:"consumer"`
}

type KafkaProducerConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
}

type KafkaConsumerConfig struct {
	Brokers []string `mapstructure:"brokers"`
	GroupID string   `mapstructure:"group_id"`
	Topic   string   `mapstructure:"topic"`
	Offset  string   `mapstructure:"offset"`
}

func NewConfiguration() *Configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Error("Config file not found", "error", err)
			panic(err)
		} else {
			slog.Error("Error reading config file", "error", err)
			panic(err)
		}
	}

	var config Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		slog.Error("Unable to decode into struct", "error", err)
		panic(err)
	}

	slog.Debug("Configuration loaded", "config", config)

	return &config
}
