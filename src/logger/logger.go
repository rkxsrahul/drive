package logger

import (
	"log"
	"os"
	"sync"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var (
	// Log is global logger
	Log *zap.Logger

	// onceInit guarantee initialize logger only once
	onceInit sync.Once
)

// constants defined for more then one occurances of string
const (
	json       string = "json"
	production string = "production"
)

// Init initializes log by input parameters
// lvl - global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
// logType - set format of logs -> json or tab space
// env - set logs environment on basis of service environment -> development or production
func Init(lvl int, logType, env string) error {
	var err error

	onceInit.Do(func() {
		// First, define our level-handling logic.
		globalLevel := zapcore.Level(lvl)

		// High-priority output should also go to standard error, and low-priority
		// output should also go to standard out.
		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= globalLevel && lvl < zapcore.ErrorLevel
		})
		// High-priority output should also go to standard error, and low-priority
		// output should also go to standard out.
		consoleDebugging := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		// default will be tab space formatting and production environment of logs
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
		if logType != json && env == production {
			consoleEncoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		} else if logType == json && env != production {
			consoleEncoder = zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
		} else if logType == json && env == production {
			consoleEncoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		}
		// Join the outputs, encoders, and level-handling functions into
		// zapcore.Cores, then tee the four cores together.
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
		)

		// From a zapcore.Core, it's easy to construct a Logger.
		Log = zap.New(core)
		zap.RedirectStdLog(Log)
		zap.ReplaceGlobals(Log)

	})

	return err
}

// AddLogging returns grpc.Server config option that turn on logging.
func AddLogging(level int, logType, env string) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, error) {
	if err := Init(level, logType, env); err != nil {
		log.Println(err)
		return nil, nil, err
	}

	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLogger(Log)

	// return unary server interceptor and stream server interceptor
	return grpc_zap.UnaryServerInterceptor(Log), grpc_zap.StreamServerInterceptor(Log), nil
}
