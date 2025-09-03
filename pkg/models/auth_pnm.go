package models

import (
	"errors"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/retawsolit/WeMeet-protocol/auth"
	"github.com/retawsolit/WeMeet-protocol/wemeet"
)

func (m *AuthModel) GeneratePNMJoinToken(c *wemeet.WeMeetTokenClaims) (string, error) {
	return auth.GenerateWeMeetJWTAccessToken(m.app.Client.ApiKey, m.app.Client.Secret, c.UserId, *m.app.Client.TokenValidity, c)
}

func (m *AuthModel) VerifyWeMeetAccessToken(token string, withTime bool) (*wemeet.WeMeetTokenClaims, error) {
	return auth.VerifyWeMeetAccessToken(m.app.Client.ApiKey, m.app.Client.Secret, token, withTime)
}

func (m *AuthModel) UnsafeClaimsWithoutVerification(token string) (*wemeet.WeMeetTokenClaims, error) {
	cl := new(wemeet.WeMeetTokenClaims)
	tk, err := jwt.ParseSigned(token, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return nil, err
	}

	err = tk.UnsafeClaimsWithoutVerification(cl)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

// RenewPNMToken we'll renew token
func (m *AuthModel) RenewPNMToken(token string, withTime bool) (string, error) {
	claims, err := m.VerifyWeMeetAccessToken(token, withTime)
	if err != nil {
		return "", err
	}

	status, err := m.natsService.GetRoomUserStatus(claims.RoomId, claims.UserId)
	if err != nil {
		return "", err
	}
	if status == "" {
		return "", errors.New("user not found")
	}

	return auth.GenerateWeMeetJWTAccessToken(m.app.Client.ApiKey, m.app.Client.Secret, claims.UserId, *m.app.Client.TokenValidity, claims)
}
