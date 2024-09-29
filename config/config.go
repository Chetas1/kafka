package config

import "github.com/spf13/viper"

type Config struct {
	Kafka         Kafka
	KafkaProducer KafkaProducer
	KafkaConsumer KafkaConsumer
}

type Kafka struct {
	Broker   string
	Username string
	Password string
}

type KafkaProducer struct {
	Topic string
}

type KafkaConsumer struct {
	Topic string
	Group string
}

func LoadConfig(path string) (Config, error) {
	var config Config

	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}

func GetConfig() (Config, error) {
	return LoadConfig("./config/config-local.yaml")
}
