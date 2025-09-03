package models

import (
	"errors"

	"github.com/retawsolit/WeMeet-protocol/wemeet"
	"github.com/retawsolit/WeMeet-server/pkg/config"
)

func (m *SpeechToTextModel) SpeechToTextTranslationServiceStart(r *wemeet.SpeechToTextTranslationReq) error {
	if !config.GetConfig().AzureCognitiveServicesSpeech.Enabled {
		return errors.New("speech service disabled")
	}

	meta, err := m.natsService.GetRoomMetadataStruct(r.RoomId)
	if err != nil {
		return err
	}
	if meta == nil {
		return errors.New("invalid nil room metadata information")
	}

	f := meta.RoomFeatures.SpeechToTextTranslationFeatures

	f.IsEnabled = r.IsEnabled
	f.AllowedSpeechLangs = r.AllowedSpeechLangs
	f.AllowedSpeechUsers = r.AllowedSpeechUsers

	f.IsEnabledTranslation = r.IsEnabledTranslation
	f.AllowedTransLangs = r.AllowedTransLangs
	f.DefaultSubtitleLang = r.DefaultSubtitleLang

	err = m.natsService.UpdateAndBroadcastRoomMetadata(r.RoomId, meta)
	if err != nil {
		return err
	}

	// send analytics
	val := wemeet.AnalyticsStatus_ANALYTICS_STATUS_STARTED.String()
	d := &wemeet.AnalyticsDataMsg{
		EventType: wemeet.AnalyticsEventType_ANALYTICS_EVENT_TYPE_ROOM,
		EventName: wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_SPEECH_SERVICE_STATUS,
		RoomId:    r.RoomId,
		HsetValue: &val,
	}
	if !f.IsEnabled {
		val = wemeet.AnalyticsStatus_ANALYTICS_STATUS_ENDED.String()
		d.EventName = wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_SPEECH_SERVICE_STATUS
		d.HsetValue = &val
	}
	m.analyticsModel.HandleEvent(d)

	return nil
}
