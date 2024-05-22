package config

import "os"

var (
	Server Ser
)

func LoadConfig(param ...string) {
	//全局的环境配置地方
	if len(param) > 0 {
		appMode = param[0]
	} else {
		appMode = os.Getenv(appModeEnv)
	}
	Server.init()
}
