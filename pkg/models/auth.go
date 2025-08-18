package models

import (
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	"github.com/retawsolit/plugnmeet-server/pkg/config"
)

type AuthModel struct {
	app         *config.AppConfig
	natsService *natsservice.NatsService
}

func NewAuthModel(app *config.AppConfig, natsService *natsservice.NatsService) *AuthModel {
	if app == nil {
		app = config.GetConfig()
	}
	if natsService == nil {
		natsService = natsservice.New(app)
	}

	return &AuthModel{
		app:         config.GetConfig(),
		natsService: natsService,
	}
}
