package log

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	{
		err := os.MkdirAll("/var/log/cf-ddns", 0666)
		if err != nil {
			log.Fatalf("Error when create log dir, details: %v", err)
		}
	}

	{
		f, err := os.OpenFile("/var/log/cf-ddns/cf-ddns.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error when opening file, details: %v", err)
		}
		f.Close()
	}

	{
		f, err := os.OpenFile("/var/log/cf-ddns/cf-ddns.err.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	cfg.OutputPaths = []string{"stdout", "/var/log/cf-ddns/cf-ddns.log"}
	cfg.ErrorOutputPaths = []string{"stdout", "/var/log/cf-ddns/cf-ddns.err.log"}
	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Error when init logger, details: %v", err)
	}

	Logger = logger.Sugar()
}