package otel_zap_logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3" // 定时任务库
	"go.uber.org/zap"           // Package zap provides fast, structured, leveled logging.
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2" // Lumberjack is a Go package for writing logs to rolling files.
)

// var zlogger *zap.Logger
type zapLogger struct {
	Logger    *zap.Logger
	lumLogger *lumberjack.Logger
}

// Go 中不仅有 channel 这种 CSP 同步机制，还有 sync.Mutex、sync.WaitGroup 等比较原始的同步原语。
// 使用它们，可以更灵活的控制数据同步和多协程并发。
var rotateCrondOnce sync.Once

// newZLogger init a zap logger;format、level等都需要配置
func newZapLogger(conf Config) zapLogger {
	// maxage default 7 days
	if config.MaxAge == 0 {
		config.MaxAge = 7
	}
	// log rolling config
	hook := lumberjack.Logger{
		Filename:   conf.File,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
		LocalTime:  true,
		Compress:   conf.Compress,
	}
	// Multi writer
	// lumberWriter and consoleWrite
	var multiWriter zapcore.WriteSyncer
	var writeSyncers []zapcore.WriteSyncer

	if conf.EnableLog {
		writeSyncers = append(writeSyncers, zapcore.AddSync(&hook)) // file log
	}

	writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout)) // console log

	if len(writeSyncers) > 0 {
		multiWriter = zapcore.NewMultiWriteSyncer(writeSyncers...)
	}

	// encoderConfig
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		// EncodeLevel:   zapcore.LowercaseLevelEncoder, // change by Himan 15Jul2022
		EncodeLevel: cEncodeLevel,
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// logLevel
	// Encoder console or json
	// enco := zapcore.NewJSONEncoder(encoderConfig) // change by Himan 15Jul2022
	enco := zapcore.NewConsoleEncoder(encoderConfig)

	var atomicLevel zap.AtomicLevel

	if conf.Debug {
		atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		atomicLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	// debug mode,use console encoder // comment by Himan 15Jul2022
	// {
	// 	_, err := os.Stat("./__debug_bin")
	// 	if err == nil {
	// 		enco = zapcore.NewConsoleEncoder(encoderConfig)
	// 	}
	// }

	// new core config
	core := zapcore.NewCore(
		enco,
		multiWriter,
		atomicLevel,
	)

	// new logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return zapLogger{
		Logger:    logger,
		lumLogger: &hook,
	}
}

// rotateCrond，多长时间一个文件这种东西要自己写？
func (zl zapLogger) rotateCrond(conf Config) {
	if conf.Rotate != "" {
		rotateCrondOnce.Do(func() {
			cron := cron.New(cron.WithSeconds())
			cron.AddFunc(conf.Rotate, func() {
				fmt.Println("rotate")
				zl.lumLogger.Rotate()
			})
			cron.Start()
		})
	}
}

// cEncodeLevel 自定义日志级别显示 // add by Himan 15Jul2022
func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
