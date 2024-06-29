package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_CfgValidate(t *testing.T) {

	cfg1 := GetDefaultConfig()
	require.NoError(t, validate.Struct(cfg1))

	cfg1 = GetDefaultConfig()
	cfg1.TemporalHostPort = "localhost`"
	require.Error(t, validate.Struct(cfg1))

	cfg1 = GetDefaultConfig()
	cfg1.KafkaBrokers = []string{}
	require.Error(t, validate.Struct(cfg1))

	cfg1 = GetDefaultConfig()
	cfg1.KafkaBrokers = []string{"1234"}
	require.Error(t, validate.Struct(cfg1))

	cfg1 = GetDefaultConfig()
	cfg1.ProducerNumber = 0
	require.Error(t, validate.Struct(cfg1))

	cfg1 = GetDefaultConfig()
	cfg1.ConsumerNumber = -1
	require.Error(t, validate.Struct(cfg1))

}

func Test_loadConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // 在测试结束后删除临时文件

	// 写入配置内容
	configContent := `
temporal_host_port: "localhost:7233"
producer_number: 2
consumer_number: 2
zookeeper_host_port: "localhost:2181"
kafka_brokers: ["localhost:32400", "localhost:32401", "localhost:32402"]
`
	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	// 调用 readInConfig 函数
	require.NoError(t, loadConfig(tmpFile.Name()))
}
