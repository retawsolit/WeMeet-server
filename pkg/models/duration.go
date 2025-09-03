package models

import (
	"github.com/retawsolit/WeMeet-server/pkg/config"
	dbservice "github.com/retawsolit/WeMeet-server/pkg/services/db"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
)

type RoomDurationModel struct {
	app         *config.AppConfig
	ds          *dbservice.DatabaseService
	rs          *redisservice.RedisService
	natsService *natsservice.NatsService
}

func NewRoomDurationModel(app *config.AppConfig, rs *redisservice.RedisService) *RoomDurationModel {
	if app == nil {
		app = config.GetConfig()
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}

	return &RoomDurationModel{
		app:         app,
		rs:          rs,
		natsService: natsservice.New(app),
	}
}
