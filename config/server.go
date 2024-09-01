package config

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	appModeEnv = "APP_MODE"
	appName    = "chat_agent" //项目名称-目前本地测试需要
	AppName    = appName

	//框架环境 stable|TestingMode|DevelopMode|LocalMode
	ReleaseMode  = "stable"
	TestingMode  = "testing"
	DevelopMode  = "develop"
	LocalMode    = "local"
	serverConfig = "server.yaml"
)

var appMode = "stable"
var configPath = "/etc/config/"

type Ser struct {
	Name      string `yaml:"name"`
	Mode      string `yaml:"mode"`
	Transport `yaml:"transport"`
}

type Transport struct {
	Http `yaml:"http"`
	Grpc `yaml:"grpc"`
}

type Http struct {
	Port int    `yaml:"port"`
	Ip   string `yaml:"ip"`
}

type Grpc struct {
	Port int    `yaml:"port"`
	Ip   string `yaml:"ip"`
}
type GrpcServer struct {
	Port int    `yaml:"port"`
	Ip   string `yaml:"ip"`
}

func (config *Ser) init() {
	//appMode = os.Getenv(appModeEnv)
	//是否走本地配置文件
	configMapMode := map[string]string{
		"local":   "1",
		"testcmd": "1",
	}
	if _, ok := configMapMode[appMode]; ok {
		workPath, err := os.Getwd()
		if err != nil {
			log.Panicf("get work path error %v", err)
		}
		findIndex := strings.Index(workPath, appName)
		if findIndex >= 0 {
			workPath = filepath.Join(workPath[:findIndex], appName)
		}
		configPath = filepath.Join(workPath, "config", appMode)
	}
	filePath := filepath.Join(configPath, serverConfig)
	yamlContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Panicf("load %s error %v", serverConfig, err)
	}
	err = yaml.Unmarshal(yamlContent, config)
	if err != nil {
		log.Panicf("unmarshal %s error %v", serverConfig, err)
	}
	switch appMode {
	case ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	case TestingMode:
		gin.SetMode(gin.TestMode)
	case DevelopMode:
		fallthrough
	case LocalMode:
		gin.SetMode(gin.DebugMode)
	default:
		log.Panicf("undefined app mode %v", config.Mode)
	}
}

// 获取环境变量
func GetAppModeEnv() string {
	return appModeEnv
}

// 获取当前app模式（stable|TestingMode|DevelopMode|LocalMode）
func GetAppMode() string {
	return appMode
}

// 是否Debug环境
func IsDebugEnv() bool {
	//appMode = os.Getenv(appModeEnv)
	debugEnv := map[string]struct{}{
		TestingMode: {},
		DevelopMode: {},
		LocalMode:   {},
	}
	res := false
	if _, ok := debugEnv[appMode]; ok {
		res = true
	}
	return res
}

// 是否开发环境
func IsDevelopEnv() bool {
	debugEnv := map[string]struct{}{
		DevelopMode: {},
		LocalMode:   {},
	}
	res := false
	if _, ok := debugEnv[appMode]; ok {
		res = true
	}
	return res
}
