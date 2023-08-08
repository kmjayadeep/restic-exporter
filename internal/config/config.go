package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Repos []ResticRepository `yaml:"repos"`
}

type ResticRepository struct {
	Name           string `yaml:"name"`
	AccessKey      string `yaml:"accessKey"`
	SecretKey      string `yaml:"secretKey"`
	Endpoint       string `yaml:"endpoint"`
	ResticPassword string `yaml:"resticPassword"`
}

func ParseConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}

	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}
