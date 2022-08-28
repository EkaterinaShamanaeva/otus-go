package config

import (
	"gopkg.in/yaml.v2" //nolint:typecheck
	"os"
)

type SchedulerConfig struct {
	Logger struct {
		Level string `yaml:"level"`
		Path  string `yaml:"path"`
	} `yaml:"logger"`
	Schedule struct {
		Period    string `yaml:"period"`
		RemindFor string `yaml:"remind_for"`
		Uri       string `yaml:"uri"`
		Queue     string `yaml:"queue"`
	} `yaml:"schedule"`
	Server struct {
		Host     string `yaml:"host"`
		HttpPort string `yaml:"http_port"`
		GrpcPort string `yaml:"grpc_port"`
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

func NewConfigScheduler() SchedulerConfig {
	return SchedulerConfig{}
}

func (config *SchedulerConfig) BuildConfigScheduler(path string) error {
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
