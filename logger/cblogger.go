package logger

import (
	"chat_agent/config"
	"context"
	"github.com/zsyu9779/myUtil/logger"
	"path"
)

// InitLog 服务启动的时候初始化业务日志
func InitCbLog() {
	logConf := config.CbLog
	collectPath := path.Join(logConf.Dir, logConf.LogCollectPath)
	msgPath := path.Join(logConf.Dir, logConf.LogMsgPath)
	resourcePath := path.Join(logConf.Dir, logConf.LogSqlPath)
	logger.InitLog(collectPath, msgPath, resourcePath)
}

// LogCollect 大数据日志收集
func LogCollect(ctx context.Context, dbname, tablename string, fields map[string]interface{}) {
	logger.LogCollect(ctx, dbname, tablename, fields)
}

// LogMsg 业务日志收集
//
//	fileName 格式为 xxx 或 xxx/xxx xxx/.../xxx, 会自动追加 _Ymd.log
func LogMsg(ctx context.Context, fileName string, args ...interface{}) {
	logger.LogMsg(ctx, fileName, args...)
}

// LogMsgNoCut 业务日志收集
//
//	fileName 格式为 xxx 或 xxx/xxx xxx/.../xxx, 自动追加 .log，但是不切割
func LogMsgNoCut(ctx context.Context, fileName string, args ...interface{}) {
	logger.LogMsgNoCut(ctx, fileName, args...)
}

// LogJson 记录json日志
//
//	fileName 格式为 xxx 或 xxx/xxx xxx/.../xxx, 会自动追加 _Ymd.log
//	fields 可为 map[string]interface{} 或 jsonBytes
func LogJson(ctx context.Context, fileName string, fields interface{}) {
	logger.LogJson(ctx, fileName, fields)
}

// LogJsonNoCut 记录json日志
//
//	fileName 格式为 xxx 或 xxx/xxx xxx/.../xxx, 自动追加 .log，但是不切割
//	fields 可为 map[string]interface{} 或 jsonBytes
func LogJsonNoCut(ctx context.Context, fileName string, fields interface{}) {
	logger.LogJsonNoCut(ctx, fileName, fields)
}
