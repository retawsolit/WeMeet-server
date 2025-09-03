package models

import (
	"context"
	"errors"

	"github.com/retawsolit/WeMeet-protocol/wemeet"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
)

func (m *BreakoutRoomModel) JoinBreakoutRoom(ctx context.Context, r *wemeet.JoinBreakoutRoomReq) (string, error) {
	status, err := m.natsService.GetRoomUserStatus(r.BreakoutRoomId, r.UserId)
	if err != nil {
		return "", err
	}
	if status == natsservice.UserStatusOnline {
		return "", errors.New("user has already been joined")
	}

	room, err := m.fetchBreakoutRoom(r.RoomId, r.BreakoutRoomId)
	if err != nil {
		return "", err
	}

	if !r.IsAdmin {
		canJoin := false
		for _, u := range room.Users {
			if u.Id == r.UserId {
				canJoin = true
			}
		}
		if !canJoin {
			return "", errors.New("you can't join in this room")
		}
	}

	p, meta, err := m.natsService.GetUserWithMetadata(r.RoomId, r.UserId)
	if err != nil {
		return "", err
	}

	req := &wemeet.GenerateTokenReq{
		RoomId: r.BreakoutRoomId,
		UserInfo: &wemeet.UserInfo{
			UserId:       r.UserId,
			Name:         p.Name,
			IsAdmin:      meta.IsAdmin,
			UserMetadata: meta,
		},
	}
	um := NewUserModel(m.app, m.ds, m.rs)
	token, err := um.GetPNMJoinToken(ctx, req)
	if err != nil {
		return "", err
	}

	return token, nil
}
