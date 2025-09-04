package models

import (
	"context"
	"errors"

	"github.com/retawsolit/wemeet-protocol/utils"
	"github.com/retawsolit/wemeet-protocol/wemeet"
)

func (m *LtiV1Model) LTIV1JoinRoom(ctx context.Context, c *wemeet.LtiClaims) (string, error) {
	res, _, _, _ := m.rm.IsRoomActive(ctx, &wemeet.IsRoomActiveReq{
		RoomId: c.RoomId,
	})

	if !res.GetIsActive() {
		_, err := m.createRoomSession(ctx, c)
		if err != nil {
			return "", errors.New(err.Error())
		}
	}

	token, err := m.joinRoom(ctx, c)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *LtiV1Model) createRoomSession(ctx context.Context, c *wemeet.LtiClaims) (*wemeet.ActiveRoomInfo, error) {
	req := utils.PrepareLTIV1RoomCreateReq(c)
	return m.rm.CreateRoom(ctx, req)
}

func (m *LtiV1Model) joinRoom(ctx context.Context, c *wemeet.LtiClaims) (string, error) {
	um := NewUserModel(m.app, m.ds, m.rs)
	token, err := um.GetPNMJoinToken(ctx, &wemeet.GenerateTokenReq{
		RoomId: c.RoomId,
		UserInfo: &wemeet.UserInfo{
			UserId:  c.UserId,
			Name:    c.Name,
			IsAdmin: c.IsAdmin,
			UserMetadata: &wemeet.UserMetadata{
				IsAdmin: c.IsAdmin,
			},
		},
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
