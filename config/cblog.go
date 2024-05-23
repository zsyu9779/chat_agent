package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v3"
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
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panicf("load %s error %v", logConfig, err)
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic("load cblog conf error " + err.Error())
	}
}
