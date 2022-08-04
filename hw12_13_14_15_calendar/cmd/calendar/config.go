package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `yaml:"logger"`
	// TODO
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		Name        string `yaml:"name"`
		Connections int    `yaml:"connections"`
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
	d := yaml.NewDecoder(f)

	// Start YAML decoding from file
	if err = d.Decode(&config); err != nil {
		return err
	}

	return nil
}

// ParseFlags parse the CLI flags and return the config path // TODO check
func ParseFlags() (string, error) {
	// String that contains configured config path
	var configPath string

	flag.StringVar(&configPath, "config", "./config.yaml", "path to config file")
	flag.Parse()

	// Make sure, that the path provided is a file, that can be read
	if err := CheckPath(configPath); err != nil {
		return "", err
	}

	return configPath, nil
}

// CheckPath makes sure, that the path provided is a file, that can be read
func CheckPath(configPath string) error {
	s, err := os.Stat(configPath)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("%s is a directory, not a file", configPath)
	}
	return nil
}
