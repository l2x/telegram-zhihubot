package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	HelpMsg = `
/s@zhihubot 关键字
	`
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
		Host            string
		SearchResultNum int `yaml:"search_result_num"`
		InlineResultNum int `yaml:"inline_result_num"`
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
	if cfg.Zhihu.SearchResultNum == 0 {
		cfg.Zhihu.SearchResultNum = 5
	}
	if cfg.Zhihu.InlineResultNum == 0 {
		cfg.Zhihu.InlineResultNum = 10
	}
	return nil
}
