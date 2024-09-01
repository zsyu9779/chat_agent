package logging

import (
	"chat_agent/logger"
	"io"

	"github.com/sirupsen/logrus"
)

func WriterLevel(lvl string) *io.PipeWriter {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		logger.Panic(err)
	}
	//logging 本包已经声明
	return logger.Logger.WriterLevel(level)
}
