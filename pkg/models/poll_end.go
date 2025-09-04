package models

import (
	"github.com/retawsolit/wemeet-protocol/wemeet"
	log "github.com/sirupsen/logrus"
)

func (m *PollModel) ClosePoll(r *wemeet.ClosePollReq) error {
	err := m.rs.ClosePoll(r)
	if err != nil {
		return err
	}

	err = m.natsService.BroadcastSystemEventToRoom(wemeet.NatsMsgServerToClientEvents_POLL_CLOSED, r.RoomId, r.PollId, nil)
	if err != nil {
		log.Errorln(err)
	}

	// send analytics
	m.analyticsModel.HandleEvent(&wemeet.AnalyticsDataMsg{
		EventType: wemeet.AnalyticsEventType_ANALYTICS_EVENT_TYPE_ROOM,
		EventName: wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_POLL_ENDED,
		RoomId:    r.RoomId,
		HsetValue: &r.PollId,
	})

	return nil
}

func (m *PollModel) CleanUpPolls(roomId string) error {
	polls, err := m.ListPolls(roomId)
	if err != nil {
		return err
	}

	var pIds []string
	for _, p := range polls {
		pIds = append(pIds, p.Id)
	}

	return m.rs.CleanUpPolls(roomId, pIds)
}
