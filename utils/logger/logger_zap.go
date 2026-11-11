package logger

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"strings"

	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vitor-franz/book-log/entities"
)

// CreateLogHandler returns a Zap-backed slog handler. Dev uses colored console
// output on stdout; prod and homo use JSON at outputPath, or stderr when empty.
// An empty environment defaults to dev. Invalid environments and output paths
// return an error and a nil handler.
func CreateLogHandler(outputPath string, environment entities.RuntimeEnvironment) (slog.Handler, error) {
	switch entities.RuntimeEnvironment(strings.ToLower(string(environment))) {
	case "", entities.EnvironmentDevelopment:
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder := zapcore.NewConsoleEncoder(encoderConfig)
		core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
		zapLogger := zap.New(core, zap.AddCaller())

		return slogzap.Option{
			Level:           slog.LevelDebug,
			Logger:          zapLogger,
			AddSource:       true,
			AttrFromContext: []func(context.Context) []slog.Attr{GetLogValuesFromContext},
		}.NewZapHandler(), nil

	case entities.EnvironmentProduction, entities.EnvironmentHomologation:
		configuration := zap.NewProductionConfig()
		// Keep an access record for every request, even during traffic bursts.
		configuration.Sampling = nil
		if outputPath != "" {
			configuration.OutputPaths = []string{outputPath}
		}

		zapLogger, err := configuration.Build()
		if err != nil {
			return nil, errors.Join(errors.New("failed to build production logger"), err)
		}

		return slogzap.Option{
			Level:           slog.LevelInfo,
			Logger:          zapLogger,
			AddSource:       true,
			AttrFromContext: []func(context.Context) []slog.Attr{GetLogValuesFromContext},
		}.NewZapHandler(), nil

	default:
		return nil, errors.New("unknown logging environment: " + string(environment))
	}
}
