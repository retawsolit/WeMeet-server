package models

import (
	"github.com/retawsolit/WeMeet-server/pkg/config"
	dbservice "github.com/retawsolit/WeMeet-server/pkg/services/db"
	livekitservice "github.com/retawsolit/WeMeet-server/pkg/services/livekit"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
)

type IngressModel struct {
	app         *config.AppConfig
	ds          *dbservice.DatabaseService
	rs          *redisservice.RedisService
	lk          *livekitservice.LivekitService
	natsService *natsservice.NatsService
}

func NewIngressModel(app *config.AppConfig, ds *dbservice.DatabaseService, rs *redisservice.RedisService, lk *livekitservice.LivekitService) *IngressModel {
	if app == nil {
		app = config.GetConfig()
	}
	if ds == nil {
		ds = dbservice.New(app.DB)
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}
	if lk == nil {
		lk = livekitservice.New(app)
	}

	return &IngressModel{
		app:         app,
		ds:          ds,
		rs:          rs,
		lk:          lk,
		natsService: natsservice.New(app),
	}
}
