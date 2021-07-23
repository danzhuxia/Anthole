package common

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type CommonStruct struct {
	Token string `json:"token" yaml:"token"`
}

type ServerStruct struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

type ClientStruct struct {
	Services []ServiceStruct `json:"services" yaml:"services"`
}

type ServiceStruct struct {
	LocalHost  string `json:"local_host" yaml:"local_host"`
	LocalPort  int    `json:"local_port" yaml:"local_port"`
	RemotePort int    `json:"remote_port" yaml:"remote_port"`
	Type       string `json:"type" yaml:"type"`
}

type AntHoleConfig struct {
	Common CommonStruct `json:"common" yaml:"common"`
	Server ServerStruct `json:"server" yaml:"server"`
	Client ClientStruct `json:"client" yaml:"client"`
}

var AntConf *AntHoleConfig

func GetConfig(configPath string) (*AntHoleConfig, error) {
	if configPath == "" {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("create config file failed: %s", err.Error())
		}

		// Search config in home directory with name ".anthole" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".anthole")
		return nil, fmt.Errorf("create config file success, check your home dictionary and filter your config: %s", home)
	}

	// Use config file from the flag.
	viper.SetConfigFile(configPath)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		return nil, fmt.Errorf("using config file: %v", viper.ConfigFileUsed())
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
