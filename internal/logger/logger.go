package logger

import (
	"cloudflare-ddns/internal/config"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logDir     string
	stdoutPath string
	stderrPath string
)

func init() {
	if config.Development {
		logDir = "./log"
	} else if config.LogDir != "" {
		logDir = config.LogDir
	} else {
		logDir = "/var/log/cf-ddns"
	}
	stdoutPath = fmt.Sprintf("%s/out.log", logDir)
	stderrPath = fmt.Sprintf("%s/err.log", logDir)

	{
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			log.Fatalf("Error when create log dir, details: %v", err)
		}
	}

	{
		f, err := os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalf("Error when opening file, details: %v", err)
		}
		f.Close()
	}

	{
		f, err := os.OpenFile(stderrPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalf("Error when opening file, details: %v", err)
		}
		f.Close()
	}
}

var Logger *zap.SugaredLogger

func InitLogger() {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encoderCfg
	cfg.OutputPaths = []string{"stdout", stdoutPath}
	cfg.ErrorOutputPaths = []string{"stdout", stderrPath}
	cfg.Sampling = nil
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Error when init logger, details: %v", err)
	}

	Logger = logger.Sugar()
}
