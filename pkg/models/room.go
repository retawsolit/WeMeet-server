package models

import (
	"github.com/retawsolit/WeMeet-server/pkg/config"
	dbservice "github.com/retawsolit/WeMeet-server/pkg/services/db"
	livekitservice "github.com/retawsolit/WeMeet-server/pkg/services/livekit"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
)

type RoomModel struct {
	app         *config.AppConfig
	ds          *dbservice.DatabaseService
	rs          *redisservice.RedisService
	lk          *livekitservice.LivekitService
	userModel   *UserModel
	natsService *natsservice.NatsService
}

func NewRoomModel(app *config.AppConfig, ds *dbservice.DatabaseService, rs *redisservice.RedisService) *RoomModel {
	if app == nil {
		app = config.GetConfig()
	}
	if ds == nil {
		ds = dbservice.New(app.DB)
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}

	return &RoomModel{
		app:         app,
		ds:          ds,
		rs:          rs,
		lk:          livekitservice.New(app),
		userModel:   NewUserModel(app, ds, rs),
		natsService: natsservice.New(app),
	}
}
