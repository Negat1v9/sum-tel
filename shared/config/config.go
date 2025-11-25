package config

import (
	"errors"

	"github.com/spf13/viper"
)

type CoreConfig struct {
	AppConfig
	WebConfig
	GrpcClientConfig
	PostgresConfig
}

type ParserServiceConfig struct {
	AppConfig
	GrpcServerConfig
	PostgresConfig
	RawMessageProducerConfig
}

type GrpcServerConfig struct {
	GPRPCPort int
}

type RawMessageProducerConfig struct {
	Topic     string
	Broker    string
	BatchSize int
}

type AppConfig struct {
	Env string
}

type WebConfig struct {
	ListenAddress string
	JwtSecret     string
	ReadTimeout   int64
	WriteTimeout  int64
}

type GrpcClientConfig struct {
	URL string
}

type PostgresConfig struct {
	DbHost     string
	DbPort     int
	DbName     string
	DbUser     string
	DbPassword string
	DbSslMode  string
}

func parseCfg(fileName string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(fileName)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func LoadCoreConfig(fileName string) (*CoreConfig, error) {
	v, err := parseCfg(fileName)
	if err != nil {
		return nil, err
	}
	var cfg *CoreConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadParserConfig(fileName string) (*ParserServiceConfig, error) {
	v, err := parseCfg(fileName)
	if err != nil {
		return nil, err
	}
	var cfg *ParserServiceConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

type NarratorConfig struct {
	AppConfig
	RawMessageConsumerCfg
	NarrationProducerCfg
	AiClientCfg
}

type RawMessageConsumerCfg struct {
	Topic      string
	Broker     string
	GroupID    string
	AutoCommit bool
}

type NarrationProducerCfg struct {
	Topic     string
	Broker    string
	BatchSize int
}

type AiClientCfg struct {
	BaseUrl      string
	Token        string
	Model        string
	SystemPrompt string
}

func LoadNarratorConfig(fileName string) (*NarratorConfig, error) {
	v, err := parseCfg(fileName)
	if err != nil {
		return nil, err
	}
	var cfg *NarratorConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
