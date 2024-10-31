package main

import (
	"flag"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger  *zap.Logger
	release bool
	log_fp  string
)

func Init() {
	flag.BoolVar(&release, "r", false, "If Release Mode or Debug Mode. Default: False")
	flag.StringVar(&log_fp, "l", "./log/all_", "Specify Log File Path Eg: ./log/logs")
	flag.Parse()

	if release {
		current_time := time.Now().Format("2006-01-02_15-04-05")

		log_fp = log_fp + current_time + ".log"

		file, err := os.OpenFile(log_fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("error occured in opening Log file\nerr: ", err)

		}
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "time" // Key for the timestamp
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), // or zapcore.NewConsoleEncoder
			zapcore.AddSync(file),
			zapcore.InfoLevel,
		)

		logger = zap.New(core)
	} else {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "time" // Key for the timestamp
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), // or zapcore.NewConsoleEncoder
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		)

		logger = zap.New(core)
	}

}

func Close() {
	logger.Sync()
}

func main() {
	Init()

	sugar := logger.Sugar()
	sugar.Info("~ PKr Server Started ~")

	Close()
}
