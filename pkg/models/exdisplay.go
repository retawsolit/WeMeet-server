package models

import (
	"errors"

	"github.com/retawsolit/WeMeet-protocol/wemeet"
	"github.com/retawsolit/WeMeet-server/pkg/config"
	dbservice "github.com/retawsolit/WeMeet-server/pkg/services/db"
	livekitservice "github.com/retawsolit/WeMeet-server/pkg/services/livekit"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
)

type ExDisplayModel struct {
	app         *config.AppConfig
	ds          *dbservice.DatabaseService
	rs          *redisservice.RedisService
	lk          *livekitservice.LivekitService
	natsService *natsservice.NatsService
}

func NewExDisplayModel(app *config.AppConfig, ds *dbservice.DatabaseService, rs *redisservice.RedisService) *ExDisplayModel {
	if app == nil {
		app = config.GetConfig()
	}
	if ds == nil {
		ds = dbservice.New(app.DB)
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}

	return &ExDisplayModel{
		app:         app,
		ds:          ds,
		rs:          rs,
		natsService: natsservice.New(app),
	}
}

func (m *ExDisplayModel) HandleTask(req *wemeet.ExternalDisplayLinkReq) error {
	switch req.Task {
	case wemeet.ExternalDisplayLinkTask_START_EXTERNAL_LINK:
		return m.start(req)
	case wemeet.ExternalDisplayLinkTask_STOP_EXTERNAL_LINK:
		return m.end(req)
	}

	return errors.New("not valid request")
}
