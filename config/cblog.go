package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

const (
	cbLogConfig = "cblog.yaml"
)

type cbLog struct {
	Dir            string `yaml:"Dir"`
	LogCollectPath string `yaml:"LogCollectPath"`
	LogMsgPath     string `yaml:"LogMsgPath"`
	LogSqlPath     string `yaml:"LogSqlPath"`
}

func (config *cbLog) init() {
	filePath := filepath.Join(configPath, cbLogConfig)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("load %s error %v", logConfig, err)
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic("load cblog conf error " + err.Error())
	}
}
