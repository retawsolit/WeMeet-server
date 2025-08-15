package factory

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/retawsolit/WeMeet-server/pkg/config"
)

func NewNatsConnection(cfg *config.AppConfig) error {
	if len(cfg.NATS.NatsUrls) == 0 {
		return fmt.Errorf("no NATS URLs provided")
	}

	nc, err := nats.Connect(cfg.NATS.NatsUrls[0])
	if err != nil {
		return fmt.Errorf("failed to connect NATS: %v", err)
	}

	cfg.NatsConn = nc
	return nil
}
