package logs

import (
	"fmt"
	"os"
	"testproject/task004/config"
	"testproject/task004/models/constants"
	"testproject/task004/utils"
	"time"

	"github.com/gin-gonic/gin"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {

	return Logrus()
}

func Logrus() gin.HandlerFunc {
	filePath := config.LogConfig.FilePath

	scr, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		utils.ExistDir(filePath)
	}
	logger := logrus.New()

	logger.Out = scr

	logger.SetLevel(logrus.DebugLevel)

	logWriter, _ := retalog.New(
		filePath+"%Y%m%d.log",
		retalog.WithMaxAge(7*24*time.Hour),
		retalog.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: constants.TimeFormat,
	})

	logger.AddHook(Hook)

	return func(context *gin.Context) {
		startTime := time.Now()
		context.Next()
		// 结束时间
		endTime := time.Now()
		stopTime := time.Since(startTime).Milliseconds()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		println(latencyTime)
		spendTime := fmt.Sprintf("%d ms", stopTime)
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}
		statusCode := context.Writer.Status()
		clientIp := context.ClientIP()
		userAgent := context.Request.UserAgent()
		dataSize := context.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		method := context.Request.Method
		path := context.Request.RequestURI

		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"status":    statusCode,
			"SpendTime": spendTime,
			"Ip":        clientIp,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})
		if len(context.Errors) > 0 {
			entry.Error(context.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
