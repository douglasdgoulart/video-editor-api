package configuration

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	LogLevel string `mapstructure:"log_level"`
	Logger   *slog.Logger
	Api      ApiConfig    `mapstructure:"api"`
	Kafka    KafkaConfig  `mapstructure:"kafka"`
	Job      JobConfig    `mapstructure:"job"`
	Ffmpeg   FfmpegConfig `mapstructure:"ffmpeg"`
}

type ApiConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
}

type JobConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type FfmpegConfig struct {
	Path string `mapstructure:"path"`
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

func NewLogger(logLevel string) *slog.Logger {
	var parsedlogLevel slog.Level
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		parsedlogLevel = slog.LevelDebug
	case "INFO":
		parsedlogLevel = slog.LevelInfo
	case "WARN":
		parsedlogLevel = slog.LevelWarn
	case "ERROR":
		parsedlogLevel = slog.LevelError
	default:
		parsedlogLevel = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{
		Level: parsedlogLevel,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func NewConfiguration() *Configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	logLevel := viper.GetString("log_level")
	logger := NewLogger(logLevel)

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

	config.Logger = logger

	slog.Debug("Configuration loaded", "config", config)

	return &config
}
