package models

import (
	"sync"

	"github.com/retawsolit/WeMeet-protocol/wemeet"
	"github.com/retawsolit/WeMeet-server/pkg/config"
	dbservice "github.com/retawsolit/WeMeet-server/pkg/services/db"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
)

type AnalyticsModel struct {
	sync.RWMutex
	data        *wemeet.AnalyticsDataMsg
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
