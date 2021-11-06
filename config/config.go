package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env-default:"localhost"`
		Port string `yaml:"port" env-default:"8825"`
	} `yaml:"server"`
}

var cfg *Config

func Get() *Config {
	return cfg
}

var once sync.Once

func Init() error {
	var err error

	once.Do(func() {
		fmt.Println("read application config")

		cfg = &Config{}
		if err = cleanenv.ReadConfig("config.yml", cfg); err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			err = fmt.Errorf("readConfig %w, help: %v", err, help)
		}
	})

	return err
}
