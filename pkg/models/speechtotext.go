package models

import (
	"strconv"
	"time"

	"github.com/retawsolit/WeMeet-server/pkg/config"
	"github.com/retawsolit/WeMeet-server/pkg/helpers"
	dbservice "github.com/retawsolit/WeMeet-server/pkg/services/db"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
	"github.com/retawsolit/wemeet-protocol/wemeet"
	log "github.com/sirupsen/logrus"
)

type SpeechToTextModel struct {
	app             *config.AppConfig
	ds              *dbservice.DatabaseService
	rs              *redisservice.RedisService
	analyticsModel  *AnalyticsModel
	webhookNotifier *helpers.WebhookNotifier
	natsService     *natsservice.NatsService
}

func NewSpeechToTextModel(app *config.AppConfig, ds *dbservice.DatabaseService, rs *redisservice.RedisService) *SpeechToTextModel {
	if app == nil {
		app = config.GetConfig()
	}
	if ds == nil {
		ds = dbservice.New(app.DB)
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}

	return &SpeechToTextModel{
		app:             app,
		ds:              ds,
		rs:              rs,
		analyticsModel:  NewAnalyticsModel(app, ds, rs),
		webhookNotifier: helpers.GetWebhookNotifier(app),
		natsService:     natsservice.New(app),
	}
}

func (m *SpeechToTextModel) sendToWebhookNotifier(rId, rSid string, userId *string, task wemeet.SpeechServiceUserStatusTasks, usage int64) {
	tk := task.String()
	n := m.webhookNotifier
	if n == nil {
		return
	}
	msg := &wemeet.CommonNotifyEvent{
		Event: &tk,
		Room: &wemeet.NotifyEventRoom{
			Sid:    &rSid,
			RoomId: &rId,
		},
		SpeechService: &wemeet.SpeechServiceEvent{
			UserId:     userId,
			TotalUsage: usage,
		},
	}
	err := n.SendWebhookEvent(msg)
	if err != nil {
		log.Errorln(err)
	}
}

func (m *SpeechToTextModel) OnAfterRoomEnded(roomId, sId string) error {
	if sId == "" {
		return nil
	}
	// we'll wait a little bit to make sure all users' requested has been received
	time.Sleep(config.WaitBeforeSpeechServicesOnAfterRoomEnded)

	hkeys, err := m.rs.SpeechToTextGetHashKeys(roomId)
	if err != nil {
		return err
	}
	for _, k := range hkeys {
		if k != "total_usage" {
			_ = m.SpeechServiceUsersUsage(roomId, sId, k, wemeet.SpeechServiceUserStatusTasks_SPEECH_TO_TEXT_SESSION_ENDED)
		}
	}

	// send by webhook
	usage, _ := m.rs.SpeechToTextGetTotalUsageByRoomId(roomId)
	if usage != "" {
		c, err := strconv.ParseInt(usage, 10, 64)
		if err == nil {
			m.sendToWebhookNotifier(roomId, sId, nil, wemeet.SpeechServiceUserStatusTasks_SPEECH_TO_TEXT_TOTAL_USAGE, c)
			// send analytics
			m.analyticsModel.HandleEvent(&wemeet.AnalyticsDataMsg{
				EventType:        wemeet.AnalyticsEventType_ANALYTICS_EVENT_TYPE_ROOM,
				EventName:        wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_SPEECH_SERVICE_TOTAL_USAGE,
				RoomId:           roomId,
				EventValueString: &usage,
			})
		}
	}

	// now clean
	err = m.rs.SpeechToTextDeleteRoom(roomId)
	if err != nil {
		return err
	}

	return nil
}
