package gophkeeper

import (
	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/container"
)

// Gophkeeper является объектом, инкапсулирующим в себе бизнес-логику сервиса по хранению секретов.
type Gophkeeper struct {
	Config    config.Config
	Container *container.Container
}

// New создает, конфигурирует и возвращает экземпляр объекта сервиса.
func New(cfg config.Config, cnt *container.Container) Gophkeeper {
	return Gophkeeper{Config: cfg, Container: cnt}
}
