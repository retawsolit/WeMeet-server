package models

import (
	"github.com/retawsolit/WeMeet-protocol/wemeet"
	"github.com/retawsolit/WeMeet-server/pkg/config"
	dbservice "github.com/retawsolit/WeMeet-server/pkg/services/db"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

type NatsModel struct {
	app            *config.AppConfig
	ds             *dbservice.DatabaseService
	rs             *redisservice.RedisService
	analytics      *AnalyticsModel
	authModel      *AuthModel
	natsService    *natsservice.NatsService
	userModel      *UserModel
	analyticsModel *AnalyticsModel
}

func NewNatsModel(app *config.AppConfig, ds *dbservice.DatabaseService, rs *redisservice.RedisService) *NatsModel {
	if app == nil {
		app = config.GetConfig()
	}
	if ds == nil {
		ds = dbservice.New(app.DB)
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}
	natsService := natsservice.New(app)

	return &NatsModel{
		app:            app,
		ds:             ds,
		rs:             rs,
		analytics:      NewAnalyticsModel(app, ds, rs),
		authModel:      NewAuthModel(app, natsService),
		natsService:    natsService,
		userModel:      NewUserModel(app, ds, rs),
		analyticsModel: NewAnalyticsModel(app, ds, rs),
	}
}

func (m *NatsModel) HandleFromClientToServerReq(roomId, userId string, req *wemeet.NatsMsgClientToServer) {
	switch req.Event {
	case wemeet.NatsMsgClientToServerEvents_REQ_RENEW_PNM_TOKEN:
		m.RenewPNMToken(roomId, userId, req.Msg)
	case wemeet.NatsMsgClientToServerEvents_REQ_INITIAL_DATA:
		m.HandleInitialData(roomId, userId)
	case wemeet.NatsMsgClientToServerEvents_REQ_JOINED_USERS_LIST:
		m.HandleSendUsersList(roomId, userId)
	case wemeet.NatsMsgClientToServerEvents_PING:
		m.HandleClientPing(roomId, userId)
	case wemeet.NatsMsgClientToServerEvents_REQ_RAISE_HAND:
		m.userModel.RaisedHand(roomId, userId, req.Msg)
	case wemeet.NatsMsgClientToServerEvents_REQ_LOWER_HAND:
		m.userModel.LowerHand(roomId, userId)
	case wemeet.NatsMsgClientToServerEvents_REQ_LOWER_OTHER_USER_HAND:
		m.userModel.LowerHand(roomId, req.Msg)
	case wemeet.NatsMsgClientToServerEvents_PUSH_ANALYTICS_DATA:
		ad := new(wemeet.AnalyticsDataMsg)
		err := protojson.Unmarshal([]byte(req.Msg), ad)
		if err != nil {
			log.Errorln(err)
			return
		}
		m.analytics.HandleEvent(ad)
	}
}
