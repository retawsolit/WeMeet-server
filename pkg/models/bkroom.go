package models

import (
	"github.com/retawsolit/WeMeet-server/pkg/config"
	dbservice "github.com/retawsolit/WeMeet-server/pkg/services/db"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
)

type BreakoutRoomModel struct {
	app         *config.AppConfig
	ds          *dbservice.DatabaseService
	rs          *redisservice.RedisService
	rm          *RoomModel
	natsService *natsservice.NatsService
}

func NewBreakoutRoomModel(app *config.AppConfig, ds *dbservice.DatabaseService, rs *redisservice.RedisService) *BreakoutRoomModel {
	if app == nil {
		app = config.GetConfig()
	}
	if ds == nil {
		ds = dbservice.New(app.DB)
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}

	return &BreakoutRoomModel{
		app:         app,
		ds:          ds,
		rs:          rs,
		rm:          NewRoomModel(app, ds, rs),
		natsService: natsservice.New(app),
	}
}
