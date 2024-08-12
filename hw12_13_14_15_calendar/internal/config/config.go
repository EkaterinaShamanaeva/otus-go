package config

import (
	"gopkg.in/yaml.v2" //nolint:typecheck
	"os"
)

type Config struct {
	Logger LoggerConf `yaml:"logger"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		SSLMode  string `yaml:"SSLMode"`
	} `yaml:"database"`
	Storage string `yaml:"storage"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}

func NewConfig() Config {
	return Config{}
}

func (config *Config) BuildConfig(path string) error {
	// Open the configuration file
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(f) // nolint:typecheck

	// Start YAML decoding from file
	if err = d.Decode(&config); err != nil {
		return err
	}

	return nil
}
