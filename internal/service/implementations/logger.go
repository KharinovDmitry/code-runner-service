package implementations

import (
	"errors"
	"log/slog"
	"os"
)

func NewLogger(env string) (*slog.Logger, error) {
	switch env {
	case "local":
		return slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			})), nil
	case "dev":
		file, err := os.OpenFile("app_logs", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		return slog.New(slog.NewTextHandler(
			file,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			})), nil
	case "prod":
		file, err := os.OpenFile("app_logs", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		return slog.New(slog.NewTextHandler(
			file,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			})), nil
	default:
		return nil, errors.New("Unknown ENV: " + env)
	}

}
