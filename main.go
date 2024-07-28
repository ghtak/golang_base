package main

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"time"
)

type Application struct {
	logger *zap.Logger
	fiber  *fiber.App
}

func newRollingFileLogCore() zapcore.Core {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(cfg)
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.DebugLevel
	})
	return zapcore.NewCore(encoder, writeSyncer, levelEnabler)
}

func newConsoleLogCore() zapcore.Core {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(cfg)
	writeSyncer := zapcore.Lock(os.Stderr)
	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zap.DebugLevel
	})
	return zapcore.NewCore(encoder, writeSyncer, levelEnabler)
}

func (application *Application) initLogger() {
	//var cores []zapcore.Core
	//cores = append(cores, newConsoleLogCore())
	//cores = append(cores, newConsoleLogCore())
	cores := []zapcore.Core{
		newRollingFileLogCore(),
		newConsoleLogCore(),
	}
	application.logger = zap.New(zapcore.NewTee(cores...))
	application.logger.Info("Init Logs",
		// Structured context as strongly typed Field values.
		zap.String("string", "stringvalue"),
		zap.Int("int", 3),
		zap.Duration("duration", time.Second),
	)
}

func (application *Application) initFiber() {
	application.fiber = fiber.New()
	application.fiber.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World")
	})
}

func (application *Application) run() {
	port := os.Getenv("GOLANG_BASE_APP_PORT")
	if port == "" {
		port = "3003"
	}
	log.Fatal(application.fiber.Listen(fmt.Sprintf(":%s", port)))
}

func main() {
	application := Application{}
	application.initLogger()
	application.initFiber()
	application.run()
}
