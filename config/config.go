package config

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// cfg default config.
var cfg = &Config{}

var validate = validator.New()

type Config struct {
	TemporalHostPort string `mapstructure:"temporal_host_port" validate:"required,hostname_port"`

	ProducerNumber int `mapstructure:"producer_number" validate:"required,gte=1"`
	ConsumerNumber int `mapstructure:"consumer_number" validate:"required,gte=1"`

	ZookeeperHostPort string   `mapstructure:"zookeeper_host_port" validate:"required,hostname_port"`
	KafkaBrokers      []string `mapstructure:"kafka_brokers" validate:"required,min=1,dive,hostname_port"`
}

func GetDefaultConfig() *Config {
	return &Config{
		TemporalHostPort:  "localhost:7233",
		ProducerNumber:    2,
		ConsumerNumber:    2,
		ZookeeperHostPort: "localhost:2181",
		KafkaBrokers:      []string{"localhost:32400", "localhost:32401", "localhost:32402"},
	}
}

func Init(filePath string) {
	cfg = GetDefaultConfig()
	_ = loadConfig(filePath)
}

func loadConfig(filePath string) error {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig() // 读取配置数据
	if err != nil {             // 处理读取配置文件的错误
		log.Printf("error reading config file, %s\n", err)
		return err
	}
	tmpCfg := &Config{}
	// log config
	fmt.Println(viper.AllSettings())
	err = viper.Unmarshal(tmpCfg)
	if err != nil {
		log.Printf("unable to decode into struct, %v\n", err)
		return err
	}
	if err := validate.Struct(tmpCfg); err != nil {
		log.Printf("config validation failed, %v\n", err)
		return err
	}
	cfg = tmpCfg
	fmt.Println("use new cfg: ", cfg)
	return nil
}

func GetConfig() *Config {
	return cfg
}
