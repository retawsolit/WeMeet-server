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

type ExMediaModel struct {
	app         *config.AppConfig
	ds          *dbservice.DatabaseService
	rs          *redisservice.RedisService
	lk          *livekitservice.LivekitService
	natsService *natsservice.NatsService
}

type updateRoomMetadataOpts struct {
	isActive *bool
	sharedBy *string
	url      *string
}

func NewExMediaModel(app *config.AppConfig, ds *dbservice.DatabaseService, rs *redisservice.RedisService) *ExMediaModel {
	if app == nil {
		app = config.GetConfig()
	}
	if ds == nil {
		ds = dbservice.New(app.DB)
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}

	return &ExMediaModel{
		app:         app,
		ds:          ds,
		rs:          rs,
		natsService: natsservice.New(app),
	}
}

func (m *ExMediaModel) HandleTask(req *wemeet.ExternalMediaPlayerReq) error {
	switch req.Task {
	case wemeet.ExternalMediaPlayerTask_START_PLAYBACK:
		return m.startPlayBack(req)
	case wemeet.ExternalMediaPlayerTask_END_PLAYBACK:
		return m.endPlayBack(req)
	}

	return errors.New("not valid request")
}
