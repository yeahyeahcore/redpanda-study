package config

type Config struct {
	HTTP    HTTP
	Service Service
}

type HTTP struct {
	Host string `env:"HTTP_HOST,default=localhost"`
	Port string `env:"HTTP_PORT,default=8080"`
}

type Service struct {
	Kafka Kafka `json:"kafka"`
}

type Kafka struct {
	Tariff  TariffKafka `json:"tariff"`
	Brokers []string    `json:"brokers"`
}

type TariffKafka struct {
	Topic      string `json:"topic"`
	GroupID    string `json:"groupId"`
	MessageKey string `json:"messageKey"`
}
