package app

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type App struct {
	Storage Storage
	Logger  logrus.Logger
}

type Storage interface {
	Create(s SystemStats) error
	Delete(id uuid.UUID) error
	FindAll() ([]SystemStats, error)
}
