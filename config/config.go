// Package config loads runtime configuration from a YAML file
// (`config/config-local.yaml` by default) and environment variables
// (via Viper's AutomaticEnv). Environment variables override file values
// using the same key path with `.` replaced by `_`
// (e.g. `KAFKA_USERNAME`, `KAFKA_PASSWORD`).
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config is the top-level service configuration.
type Config struct {
	Kafka         Kafka
	KafkaProducer KafkaProducer
	KafkaConsumer KafkaConsumer
}

// Kafka holds broker connection + SASL credentials.
//
// Username/Password should come from env vars in production
// (KAFKA_USERNAME / KAFKA_PASSWORD) — never commit real credentials
// to the YAML file.
type Kafka struct {
	Broker   string
	Username string
	Password string
}

// KafkaProducer holds producer-side configuration.
type KafkaProducer struct {
	Topic string
}

// KafkaConsumer holds consumer-side configuration.
type KafkaConsumer struct {
	Topic string
	Group string
}

// LoadConfig reads the YAML file at `path`, applies env overrides,
// and returns the parsed Config. Wrapped errors are returned for missing
// files or unmarshal failures.
func LoadConfig(path string) (Config, error) {
	var cfg Config

	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("read config %s: %w", path, err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}
	return cfg, nil
}

// GetConfig is a convenience wrapper that loads from `./config/config-local.yaml`.
// Production callers should use LoadConfig directly with an explicit path.
func GetConfig() (Config, error) {
	return LoadConfig("./config/config-local.yaml")
}
