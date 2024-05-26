package app

import (
	"code-runner-service/config"
	"code-runner-service/internal/service"
)

func Run(cfg *config.Config) error {
	services := service.NewServiceManager()
	if err := services.Init(cfg.Env); err != nil {
		return err
	}

	services.Logger.Info("app started")

	return nil
}
