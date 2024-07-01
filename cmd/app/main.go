package main

import (
	"code-runner-service/config"
	"code-runner-service/internal/app"
	"flag"
	"os/exec"
)

// @title Contest Service API
// @version 1.0
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "path", "", "path to config file")
	flag.Parse()
	if cfgPath == "" {
		panic("config file path is required")
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic("read config error: " + err.Error())
	}

	err = buildExecutors()
	if err != nil {
		panic("build executors error: " + err.Error())
	}

	err = app.Run(cfg)
	if err != nil {
		panic(err.Error())
	}
}

func buildExecutors() error {
	if err := exec.Command("docker", "build", "-t", "unpivileged_run", "unprivilegedRun").Run(); err != nil {
		return err
	}
	if err := exec.Command("docker", "build", "-t", "python_executor", "python").Run(); err != nil {
		return err
	}
	return nil
}
