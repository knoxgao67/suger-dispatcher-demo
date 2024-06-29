package config

var cfg *Config

type Config struct {
	HostPort string

	ProducerNumber, ConsumerNumber int

	ZookeeperHostPort string
	KafkaBrokers      []string
}

func Init() {
	cfg = &Config{
		HostPort: "localhost:7233",

		ProducerNumber: 2,
		ConsumerNumber: 2,

		//ZookeeperHostPort: "kafka-zookeeper-headless:2181",
		ZookeeperHostPort: "localhost:2181",

		KafkaBrokers: []string{"localhost:32400", "localhost:32401", "localhost:32402"},
	}
}

func GetConfig() *Config {
	return cfg
}
