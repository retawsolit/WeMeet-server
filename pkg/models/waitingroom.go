package models

import (
	"github.com/retawsolit/WeMeet-server/pkg/config"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
)

type WaitingRoomModel struct {
	app         *config.AppConfig
	rs          *redisservice.RedisService
	natsService *natsservice.NatsService
}

func NewWaitingRoomModel(app *config.AppConfig, rs *redisservice.RedisService) *WaitingRoomModel {
	if app == nil {
		app = config.GetConfig()
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}

	return &WaitingRoomModel{
		app:         app,
		rs:          rs,
		natsService: natsservice.New(app),
	}
}
