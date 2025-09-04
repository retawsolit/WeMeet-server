package natsservice

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/retawsolit/wemeet-protocol/wemeet"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func (s *NatsService) BroadcastSystemEventToRoom(event wemeet.NatsMsgServerToClientEvents, roomId string, data interface{}, toUserId *string) error {
	var msg string
	var err error

	switch v := data.(type) {
	case int:
	case float64:
		msg = fmt.Sprintf("%v", v)
	case []byte:
		msg = string(v)
	case string:
		msg = v
	case proto.Message:
		msg, err = s.MarshalToProtoJson(v)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid data type")
	}

	payload := wemeet.NatsMsgServerToClient{
		Id:    uuid.NewString(),
		Event: event,
		Msg:   msg,
	}
	message, err := proto.Marshal(&payload)
	if err != nil {
		return err
	}

	sub := fmt.Sprintf("%s:%s.system", roomId, s.app.NatsInfo.Subjects.SystemPublic)
	if toUserId != nil {
		sub = fmt.Sprintf("%s:%s.%s.system", roomId, s.app.NatsInfo.Subjects.SystemPrivate, *toUserId)
	}

	_, err = s.js.Publish(s.ctx, sub, message)
	if err != nil {
		return err
	}

	return nil
}

func (s *NatsService) BroadcastSystemEventToEveryoneExceptUserId(event wemeet.NatsMsgServerToClientEvents, roomId string, data interface{}, exceptUserId string) error {
	ids, err := s.GetOnlineUsersId(roomId)
	if err != nil {
		return err
	}
	if ids == nil || len(ids) == 0 {
		return errors.New("no online user found")
	}

	for _, id := range ids {
		if id != exceptUserId {
			go func(id string) {
				err := s.BroadcastSystemEventToRoom(event, roomId, data, &id)
				if err != nil {
					log.Errorln(err)
				}
			}(id)
		}
	}

	return nil
}

func (s *NatsService) BroadcastSystemNotificationToRoom(roomId, msg string, msgType wemeet.NatsSystemNotificationTypes, withSound bool, userId *string) error {
	data := &wemeet.NatsSystemNotification{
		Id:        uuid.NewString(),
		Type:      msgType,
		Msg:       msg,
		SentAt:    time.Now().UnixMilli(),
		WithSound: withSound,
	}

	marshal, err := protoJsonOpts.Marshal(data)
	if err != nil {
		return err
	}

	return s.BroadcastSystemEventToRoom(wemeet.NatsMsgServerToClientEvents_SYSTEM_NOTIFICATION, roomId, marshal, userId)
}

func (s *NatsService) NotifyInfoMsg(roomId, msg string, withSound bool, userId *string) error {
	return s.BroadcastSystemNotificationToRoom(roomId, msg, wemeet.NatsSystemNotificationTypes_NATS_SYSTEM_NOTIFICATION_INFO, withSound, userId)
}

func (s *NatsService) NotifyWarningMsg(roomId, msg string, withSound bool, userId *string) error {
	return s.BroadcastSystemNotificationToRoom(roomId, msg, wemeet.NatsSystemNotificationTypes_NATS_SYSTEM_NOTIFICATION_WARNING, withSound, userId)
}

func (s *NatsService) NotifyErrorMsg(roomId, msg string, userId *string) error {
	return s.BroadcastSystemNotificationToRoom(roomId, msg, wemeet.NatsSystemNotificationTypes_NATS_SYSTEM_NOTIFICATION_ERROR, true, userId)
}
