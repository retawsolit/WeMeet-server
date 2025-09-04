package factory

import (
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/retawsolit/WeMeet-server/pkg/config"
	"github.com/retawsolit/wemeet-protocol/utils"
)

func NewNatsConnection(appCnf *config.AppConfig) error {
	info := appCnf.NatsInfo
	var opt nats.Option
	var err error

	if info.Nkey != nil {
		opt, err = utils.NkeyOptionFromSeedText(*info.Nkey)
		if err != nil {
			return err
		}
	} else {
		opt = nats.UserInfo(info.User, info.Password)
	}

	nc, err := nats.Connect(strings.Join(info.NatsUrls, ","), opt)
	if err != nil {
		return err
	}
	appCnf.NatsConn = nc

	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}
	appCnf.JetStream = js

	return nil
}
