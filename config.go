package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AuthConfig struct {
	UserName string  `yaml:"username"`
	Password string  `yaml:"password"`
}

type TlsConfig struct {
	Enable  bool      `yaml:"enable"`
	CA      string    `yaml:"ca"`
	Cert    string    `yaml:"cert"`
	Key     string    `yaml:"key"`
}

type LocalConfig struct {
	Listen  string      `yaml:"listen"`
	Timeout int         `yaml:"timeout"`
	Mode    string      `yaml:"mode"`   // local、auto、proxy
	Auths []AuthConfig  `yaml:"auth"`
	Tls     TlsConfig   `yaml:"tls"`
}

type RemoteConfig struct {
	Address string       `yaml:"address"`
	Timeout int          `yaml:"timeout"`
	Auth    AuthConfig   `yaml:"auth"`
	Tls     TlsConfig    `yaml:"tls"`
}

type LogConfig struct {
	Path string   `yaml:"path"`
	FileSize int  `yaml:"filesize"`
	FileNum  int  `yaml:"filenumber"`
}

type Config struct {
	Log     LogConfig       `yaml:"log"`
	Local   LocalConfig     `yaml:"local"`
	Remote []RemoteConfig   `yaml:"remote"`
}


func LoadConfig(filename string) (*Config, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	config.Remote = make([]RemoteConfig, 0)
	err = yaml.Unmarshal(body, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}