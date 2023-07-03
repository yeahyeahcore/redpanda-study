package config

type Config struct {
	HTTP    ConfigurationHTTP
	Service ServiceConfiguration
}

type ConfigurationHTTP struct {
	Host string `env:"HTTP_HOST,default=localhost"`
	Port string `env:"HTTP_PORT,default=8080"`
}

type ServiceConfiguration struct {
	Kafka KafkaConfiguration `json:"kafka"`
}

type KafkaConfiguration struct {
	Tariff TariffKafkaConfiguration `json:"tariff"`
}

type TariffKafkaConfiguration struct {
	Topic   string   `json:"topic"`
	GroupID string   `json:"groupId"`
	Brokers []string `json:"brokers"`
}
