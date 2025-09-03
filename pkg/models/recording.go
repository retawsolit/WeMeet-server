package models

import (
	"github.com/retawsolit/WeMeet-protocol/wemeet"
	"github.com/retawsolit/WeMeet-server/pkg/config"
	"github.com/retawsolit/WeMeet-server/pkg/dbmodels"
	"github.com/retawsolit/WeMeet-server/pkg/helpers"
	dbservice "github.com/retawsolit/WeMeet-server/pkg/services/db"
	livekitservice "github.com/retawsolit/WeMeet-server/pkg/services/livekit"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	redisservice "github.com/retawsolit/WeMeet-server/pkg/services/redis"
	log "github.com/sirupsen/logrus"
)

type RecordingModel struct {
	app             *config.AppConfig
	ds              *dbservice.DatabaseService
	rs              *redisservice.RedisService
	lk              *livekitservice.LivekitService
	analyticsModel  *AnalyticsModel
	webhookNotifier *helpers.WebhookNotifier
	natsService     *natsservice.NatsService
}

func NewRecordingModel(app *config.AppConfig, ds *dbservice.DatabaseService, rs *redisservice.RedisService) *RecordingModel {
	if app == nil {
		app = config.GetConfig()
	}
	if ds == nil {
		ds = dbservice.New(app.DB)
	}
	if rs == nil {
		rs = redisservice.New(app.RDS)
	}

	return &RecordingModel{
		app:             app,
		ds:              ds,
		rs:              rs,
		analyticsModel:  NewAnalyticsModel(app, ds, rs),
		webhookNotifier: helpers.GetWebhookNotifier(app),
		natsService:     natsservice.New(app),
	}
}

func (m *RecordingModel) HandleRecorderResp(r *wemeet.RecorderToWeMeet, roomInfo *dbmodels.RoomInfo) {
	switch r.Task {
	case wemeet.RecordingTasks_START_RECORDING:
		m.recordingStarted(r)
		go m.sendToWebhookNotifier(r)

	case wemeet.RecordingTasks_END_RECORDING:
		m.recordingEnded(r)
		go m.sendToWebhookNotifier(r)

	case wemeet.RecordingTasks_START_RTMP:
		m.rtmpStarted(r)
		go m.sendToWebhookNotifier(r)

	case wemeet.RecordingTasks_END_RTMP:
		m.rtmpEnded(r)
		go m.sendToWebhookNotifier(r)

	case wemeet.RecordingTasks_RECORDING_PROCEEDED:
		creation, err := m.addRecordingInfoToDB(r, roomInfo.CreationTime)
		if err != nil {
			log.Errorln(err)
		}
		// keep record of this file
		m.addRecordingInfoFile(r, creation, roomInfo)
		go m.sendToWebhookNotifier(r)
	}
}

func (m *RecordingModel) sendToWebhookNotifier(r *wemeet.RecorderToWeMeet) {
	tk := r.Task.String()
	n := m.webhookNotifier
	if n != nil {
		msg := &wemeet.CommonNotifyEvent{
			Event: &tk,
			Room: &wemeet.NotifyEventRoom{
				Sid:    &r.RoomSid,
				RoomId: &r.RoomId,
			},
			RecordingInfo: &wemeet.RecordingInfoEvent{
				RecordId:    r.RecordingId,
				RecorderId:  r.RecorderId,
				RecorderMsg: r.Msg,
				FilePath:    &r.FilePath,
				FileSize:    &r.FileSize,
			},
		}
		if r.Task == wemeet.RecordingTasks_RECORDING_PROCEEDED {
			// this process may take longer time & webhook url may clean up
			// so, here we'll use ForceToPutInQueue method to retrieve url from mysql table
			n.ForceToPutInQueue(msg)
		} else {
			err := n.SendWebhookEvent(msg)
			if err != nil {
				log.Errorln(err)
			}
		}
	}

	// send analytics
	var val string
	data := &wemeet.AnalyticsDataMsg{
		EventType: wemeet.AnalyticsEventType_ANALYTICS_EVENT_TYPE_ROOM,
		RoomId:    r.RoomId,
	}

	switch r.Task {
	case wemeet.RecordingTasks_START_RECORDING:
		data.EventName = wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_RECORDING_STATUS
		val = wemeet.AnalyticsStatus_ANALYTICS_STATUS_STARTED.String() + ":" + r.RecorderId
	case wemeet.RecordingTasks_END_RECORDING:
		data.EventName = wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_RECORDING_STATUS
		val = wemeet.AnalyticsStatus_ANALYTICS_STATUS_ENDED.String()
	case wemeet.RecordingTasks_START_RTMP:
		data.EventName = wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_RTMP_STATUS
		val = wemeet.AnalyticsStatus_ANALYTICS_STATUS_STARTED.String() + ":" + r.RecorderId
	case wemeet.RecordingTasks_END_RTMP:
		data.EventName = wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_RTMP_STATUS
		val = wemeet.AnalyticsStatus_ANALYTICS_STATUS_ENDED.String()
	}
	data.HsetValue = &val
	m.analyticsModel.HandleEvent(data)
}
