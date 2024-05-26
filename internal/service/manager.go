package service

import (
	"code-runner-service/internal/domain/service"
	"code-runner-service/internal/service/implementations"
)

type Manager struct {
	Logger          service.Logger
	ExecutorFactory service.ExecutorFactory
	TestRunner      service.TestRunner
}

func NewServiceManager() *Manager {
	return &Manager{}
}

func (m *Manager) Init(env string) error {
	log, err := implementations.NewLogger(env)
	if err != nil {
		return err
	}
	m.Logger = log

	factory := implementations.NewExecutorFactory()

	testRunner := implementations.NewTestRunnerService(factory)
	m.TestRunner = testRunner

	return nil
}
