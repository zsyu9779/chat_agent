package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"path/filepath"
)

const (
	logConfig = "log.yaml"
)

type LogStruct struct {
	Prefix   string `json:"prefix" yaml:"prefix"`
	FilePath string `json:"file_path" yaml:"file_path"`
}

func (config *LogStruct) Init() {
	filePath := filepath.Join(configPath, logConfig)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("load %s error %v", logConfig, err)
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Panicf("unmarshal %s error %v", logConfig, err)
	}
	if config.Prefix == "" {
		config.Prefix = Server.Name
	}
}
