package models

import (
	"errors"

	"github.com/retawsolit/WeMeet-protocol/wemeet"
	"github.com/retawsolit/WeMeet-server/pkg/config"
	dbservice "github.com/retawsolit/WeMeet-server/pkg/services/db"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
)

type PollModel struct {
	app            *config.AppConfig
	ds             *dbservice.DatabaseService
	rs             *redisservice.RedisService
	analyticsModel *AnalyticsModel
	natsService    *natsservice.NatsService
}

func NewPollModel(app *config.AppConfig, ds *dbservice.DatabaseService, rs *redisservice.RedisService) *PollModel {
	if app == nil {
		app = config.GetConfig()
	}
	if ds == nil {
		ds = dbservice.New(app.DB)
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}

	return &PollModel{
		app:            app,
		ds:             ds,
		rs:             rs,
		analyticsModel: NewAnalyticsModel(app, ds, rs),
		natsService:    natsservice.New(app),
	}
}

func (m *PollModel) ManageActivation(req *wemeet.ActivatePollsReq) error {
	roomMeta, err := m.natsService.GetRoomMetadataStruct(req.GetRoomId())
	if err != nil {
		return err
	}
	if roomMeta == nil {
		return errors.New("invalid nil room metadata information")
	}

	roomMeta.RoomFeatures.PollsFeatures.IsActive = req.GetIsActive()
	return m.natsService.UpdateAndBroadcastRoomMetadata(req.GetRoomId(), roomMeta)
}
