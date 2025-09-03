package models

import (
	"fmt"
	"github.com/mynaparrot/plugnmeet-protocol/plugnmeet"
	natsservice "github.com/mynaparrot/plugnmeet-server/pkg/services/nats"
	log "github.com/sirupsen/logrus"
	"time"
)

func (m *NatsModel) OnAfterUserJoined(roomId, userId string) {
	status, err := m.natsService.GetRoomUserStatus(roomId, userId)
	if err != nil {
		log.Errorln(fmt.Sprintf("error GetRoomUserStatus to %s; roomId: %s; msg: %s", roomId, userId, err))
		return
	}
	if status == natsservice.UserStatusOnline {
		// no need
		return
	}

	err = m.natsService.UpdateUserStatus(roomId, userId, natsservice.UserStatusOnline)
	if err != nil {
		log.Warnln(fmt.Sprintf("Error updating user status: %s for %s; roomId: %s; msg: %s", natsservice.UserStatusOnline, userId, roomId, err.Error()))
	}

	if userInfo, err := m.natsService.GetUserInfo(roomId, userId); err == nil && userInfo != nil {
		// broadcast this user to everyone
		err := m.natsService.BroadcastSystemEventToEveryoneExceptUserId(plugnmeet.NatsMsgServerToClientEvents_USER_JOINED, roomId, userInfo, userId)
		if err != nil {
			log.Warnln(err)
		}
		now := fmt.Sprintf("%d", time.Now().UnixMilli())
		m.analyticsModel.HandleEvent(&plugnmeet.AnalyticsDataMsg{
			EventType: plugnmeet.AnalyticsEventType_ANALYTICS_EVENT_TYPE_ROOM,
			EventName: plugnmeet.AnalyticsEvents_ANALYTICS_EVENT_USER_JOINED,
			RoomId:    roomId,
			UserId:    &userId,
			UserName:  &userInfo.Name,
			ExtraData: &userInfo.Metadata,
			HsetValue: &now,
		})
	}
}

// OnAfterUserDisconnected should be run in separate goroutine
// we'll wait for 5 seconds before declare user as offline
// but will broadcast as disconnected
func (m *NatsModel) OnAfterUserDisconnected(roomId, userId string) {
	// now change the user's status
	_ = m.natsService.UpdateUserStatus(roomId, userId, natsservice.UserStatusDisconnected)

	userInfo, _ := m.natsService.GetUserInfo(roomId, userId)
	if userInfo == nil {
		// if we do not get data, then we'll just update analytics
		m.updateUserLeftAnalytics(roomId, userId)
		return
	}

	_ = m.natsService.BroadcastSystemEventToEveryoneExceptUserId(plugnmeet.NatsMsgServerToClientEvents_USER_DISCONNECTED, roomId, userInfo, userId)

	// Schedule to check for offline status after 5 seconds without blocking.
	time.AfterFunc(5*time.Second, func() {
		m.handleOfflineStatus(roomId, userId, userInfo)
	})
}

// handleOfflineStatus will check if the user is still disconnected and mark them as offline.
func (m *NatsModel) handleOfflineStatus(roomId, userId string, userInfo *plugnmeet.NatsKvUserInfo) {
	status, err := m.natsService.GetRoomUserStatus(roomId, userId)
	if err == nil && status == natsservice.UserStatusOnline {
		// User reconnected, do nothing.
		return
	}

	err = m.natsService.UpdateUserStatus(roomId, userId, natsservice.UserStatusOffline)
	if err != nil {
		log.Warnln(fmt.Sprintf("Error updating user status: %s for %s; roomId: %s; msg: %s", natsservice.UserStatusOffline, userId, roomId, err.Error()))
	}

	// analytics
	m.updateUserLeftAnalytics(roomId, userId)

	// now broadcast to everyone
	_ = m.natsService.BroadcastSystemEventToEveryoneExceptUserId(plugnmeet.NatsMsgServerToClientEvents_USER_OFFLINE, roomId, userInfo, userId)

	// Schedule consumer deletion after 30 seconds.
	time.AfterFunc(30*time.Second, func() {
		m.cleanupUserConsumer(roomId, userId)
	})
}

// cleanupUserConsumer will delete the NATS consumer for the user if they haven't reconnected.
func (m *NatsModel) cleanupUserConsumer(roomId, userId string) {
	status, err := m.natsService.GetRoomUserStatus(roomId, userId)
	if err == nil && status == natsservice.UserStatusOnline {
		// User reconnected, do not delete consumer.
		return
	}

	// do not need to delete the user as user may come to online again
	// when the session is ended, we'll do proper clean up.
	// delete consumer only
	m.natsService.DeleteConsumer(roomId, userId)
}

func (m *NatsModel) updateUserLeftAnalytics(roomId, userId string) {
	now := fmt.Sprintf("%d", time.Now().UnixMilli())
	m.analyticsModel.HandleEvent(&plugnmeet.AnalyticsDataMsg{
		EventType: plugnmeet.AnalyticsEventType_ANALYTICS_EVENT_TYPE_USER,
		EventName: plugnmeet.AnalyticsEvents_ANALYTICS_EVENT_USER_LEFT,
		RoomId:    roomId,
		UserId:    &userId,
		HsetValue: &now,
	})
}
