package natsservice

import (
	"errors"

	"github.com/retawsolit/wemeet-protocol/wemeet"
	log "github.com/sirupsen/logrus"
)

func (s *NatsService) BroadcastUserMetadata(roomId string, userId string, metadata, toUser *string) error {
	if metadata == nil {
		result, err := s.GetUserInfo(roomId, userId)
		if err != nil {
			return err
		}

		if result == nil {
			return errors.New("user not found")
		}
		metadata = &result.Metadata
	}

	data := &wemeet.NatsUserMetadataUpdate{
		Metadata: *metadata,
		UserId:   userId,
	}

	return s.BroadcastSystemEventToRoom(wemeet.NatsMsgServerToClientEvents_USER_METADATA_UPDATE, roomId, data, toUser)
}

// UpdateAndBroadcastUserMetadata will update metadata & broadcast to everyone
func (s *NatsService) UpdateAndBroadcastUserMetadata(roomId, userId string, meta interface{}, toUserId *string) error {
	if meta == nil {
		return errors.New("metadata cannot be nil")
	}

	mt, err := s.UpdateUserMetadata(roomId, userId, meta)
	if err != nil {
		return err
	}
	return s.BroadcastUserMetadata(roomId, userId, &mt, toUserId)
}

func (s *NatsService) BroadcastUserInfoToRoom(event wemeet.NatsMsgServerToClientEvents, roomId, userId string, userInfo *wemeet.NatsKvUserInfo) {
	if userInfo == nil {
		info, err := s.GetUserInfo(roomId, userId)
		if err != nil {
			return
		}
		if info == nil {
			return
		}
	}

	err := s.BroadcastSystemEventToRoom(event, roomId, userInfo, nil)
	if err != nil {
		log.Warnln(err)
	}
}
