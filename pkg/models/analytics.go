package models

import (
	"sync"

	"github.com/retawsolit/WeMeet-server/pkg/config"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	"github.com/retawsolit/plugnmeet-protocol/plugnmeet"
)

type AnalyticsModel struct {
	sync.RWMutex
	data        *plugnmeet.AnalyticsDataMsg
	app         *config.AppConfig
	ds          *dbservice.DatabaseService
	rs          *redisservice.RedisService
	natsService *natsservice.NatsService
}

func NewAnalyticsModel(app *config.AppConfig, ds *dbservice.DatabaseService, rs *redisservice.RedisService) *AnalyticsModel {
	if app == nil {
		app = config.GetConfig()
	}
	if ds == nil {
		ds = dbservice.New(app.DB)
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}

	return &AnalyticsModel{
		app:         config.GetConfig(),
		ds:          ds,
		rs:          rs,
		natsService: natsservice.New(app),
	}
}
