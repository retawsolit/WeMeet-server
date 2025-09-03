package models

import (
	"github.com/retawsolit/WeMeet-protocol/webhook"
)

func (m *AuthModel) ValidateLivekitWebhookToken(body []byte, token string) (bool, error) {
	return webhook.VerifyRequest(body, m.app.LivekitInfo.ApiKey, m.app.LivekitInfo.Secret, token)
}
