package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RefreshMinutes uint               `yaml:"refreshMinutes"`
	Repos          []ResticRepository `yaml:"repos"`
}

type ResticRepository struct {
	Name           string `yaml:"name"`
	AccessKey      string `yaml:"accessKey"`
	SecretKey      string `yaml:"secretKey"`
	Endpoint       string `yaml:"endpoint"`
	ResticPassword string `yaml:"resticPassword"`
	S3SizeLimit    uint   `yaml:"s3SizeLimit"`
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

	if c.RefreshMinutes == 0 {
		c.RefreshMinutes = 10
	}

	return c, nil
}
