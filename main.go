package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ButterHost69/PKr-server/db"
	"github.com/ButterHost69/PKr-server/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	router *gin.Engine
)

var (
	// Flag Variables
	RELEASE bool
	LOG_FP  string
	IPADDR  string
)

func Init() {
	flag.BoolVar(&RELEASE, "r", false, "If Release Mode or Debug Mode. Default: False")
	flag.StringVar(&LOG_FP, "l", "./log/events_", "Specify Log File Path Eg: ./log/logs")
	flag.StringVar(&IPADDR, "ipaddr", "localhost:9069", "Specify Address to Run Server")
	flag.Parse()

	if err := db.InitSQLiteDatabase(); err != nil {
		log.Fatal("error Could not start the Database.\nError: ", err)
	}

	if RELEASE {
		// Set the Logger
		current_time := time.Now().Format("2006-01-02_15-04-05")

		LOG_FP = LOG_FP + current_time + ".log"

		file, err := os.OpenFile(LOG_FP, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("error occured in opening Log file\nerr: ", err)

		}
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "time" // Key for the timestamp
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			zapcore.InfoLevel,
		)

		logger = zap.New(core)

		// Set the Gin Server
		gin.SetMode(gin.ReleaseMode)

		// TODO: [ ] Allow TLS Support, make it an argument option
		router = gin.New()
		router.Use(gin.LoggerWithWriter(zap.NewStdLog(logger).Writer()))

	} else {
		// Logger
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "time"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		)

		logger = zap.New(core)

		// Gin Router
		router = gin.Default()
	}

}

func Close() {
	logger.Sync()
	db.CloseSQLiteDatabase()
}

func main() {
	Init()

	sugar := logger.Sugar()
	sugar.Info("~ PKr Server Started ~")

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello World ... ")
	})

	router.POST("/register/user", func(ctx *gin.Context) {
		handler.RegisterUser(ctx, sugar)
	})

	router.POST("/register/workspace", func(ctx *gin.Context) {
		handler.RegisterWorkspace(ctx, sugar)
	})

	// TODO: [ ] Send 
	router.POST("/update/me", func(ctx *gin.Context) {
		handler.UserIPCheck(ctx, sugar)
	})

	if err := router.Run(IPADDR); err != nil {
		log.Fatal("error Occured in Starting Gin Server...Error: ", err)
	}

	Close()
}
