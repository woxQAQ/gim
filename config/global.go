package config

import (
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"os"
	"time"
)

type Loader interface {
	Load(name string)
}

var DB *gorm.DB

var Cfg *ini.File

var RunMode string

var (
	HttpPort     int
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	JwtSecret    string
	RecallTimes  int
)

var (
	DBName string
	DBUser string
	DBPwd  string
	DBHost string
)

var DateTemp = "2006-01-01"

var KafkaConfigPath string
var TransferConfigPath string
var GatewayConfigPath string

type defaultConfigFilesPath struct {
	KafkaConfig    string `yaml:"kafka_config"`
	TransferConfig string `yaml:"transfer_config"`
	GatewayConfig  string `yaml:"gateway_config"`
}

func init() {
	data, err := os.ReadFile("../global.yaml")
	if err != nil {
		panic(err)
	}

	var config defaultConfigFilesPath
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	KafkaConfigPath = config.KafkaConfig
	TransferConfigPath = config.TransferConfig
	GatewayConfigPath = config.GatewayConfig
}
