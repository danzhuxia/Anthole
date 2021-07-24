package common

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type CommonStruct struct {
	Token string `yaml:"token"`
}

type ServerStruct struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ClientStruct struct {
	Services []ServiceStruct `yaml:"services"`
}

type ServiceStruct struct {
	LocalHost  string `yaml:"local_host"`
	LocalPort  int    `yaml:"local_port"`
	RemotePort int    `yaml:"remote_port"`
	Type       string `yaml:"type"`
}

type AntHoleConfig struct {
	Common CommonStruct `yaml:"common"`
	Server ServerStruct `yaml:"server"`
	Client ClientStruct `yaml:"client"`
}

var AntConf = new(AntHoleConfig)

func GetConfig(configPath string) (*AntHoleConfig, error) {
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("create config file failed: %s", err.Error())
		}

		// _, err = os.Create(home + "/.anthole.yaml")

		return nil, fmt.Errorf("create config file success, check your home dictionary and filter your config: %s", home)

	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("using this config filepath: %v failed: %s", viper.ConfigFileUsed(), err.Error())
	}
	AntConf.Common.Token = viper.GetStringMap("common")["token"].(string)
	AntConf.Server.Host = viper.GetStringMap("server")["host"].(string)
	AntConf.Server.Port = viper.GetStringMap("server")["port"].(int)
	services := viper.GetStringMap("client")["services"].([]interface{})
	for _, service := range services {
		serTmp := ServiceStruct{
			LocalHost:  service.(map[interface{}]interface{})["local_host"].(string),
			LocalPort:  service.(map[interface{}]interface{})["local_port"].(int),
			RemotePort: service.(map[interface{}]interface{})["remote_port"].(int),
			Type:       service.(map[interface{}]interface{})["type"].(string),
		}
		AntConf.Client.Services = append(AntConf.Client.Services, serTmp)
	}
	return AntConf, nil
}
