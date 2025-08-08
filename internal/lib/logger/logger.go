package logger

import (
	"fmt"
	"log/slog"
	"order-service/internal/config"
	"os"
)

func InitLogger(env string) (*slog.Logger, error) {
	switch env {
	case "local", "test":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})), nil
	case "dev", "prod":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})), nil
	default:
		return nil, fmt.Errorf("%s is invalid env %w", env, config.ErrCfgInvalid)
	}

}
