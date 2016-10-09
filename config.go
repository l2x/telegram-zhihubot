package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var cfg *Config

type Config struct {
	Bot struct {
		Token string
		Debug bool
	}

	HTTP struct {
		Host       string
		Port       string
		PublicKey  string
		PrivateKey string
	}

	Zhihu struct {
		Host string
	}
}

func initConfig(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	cfg = &Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	return nil
}
