package conf

import (
	"encoding/json"
	"os"
)

type DBConfiguration struct {
	Host   string
	Port   string
	User   string
	Passwd string
	DbName string
}

type MqttConfiguration struct {
	Addr     string
	ClientID string
	TopicPub string
	TopicSub string
}

type ServerConfiguration struct {
	BindPort          string
	ReadLimit         uint16
	WriteLimit        uint16
	ConnTimeout       uint16
	ConnCheckInterval uint16
	ServerStatistics  uint16
}

type Configuration struct {
	Server *ServerConfiguration
	DB     *DBConfiguration
	Mqtt   *MqttConfiguration
}

var G_conf *Configuration

func ReadConfig(confpath string) (*Configuration, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	G_conf = &config

	return &config, err
}

func GetConf() *Configuration {
	return G_conf
}
