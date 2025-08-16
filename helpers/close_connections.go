package helpers

import (
	"github.com/retawsolit/WeMeet-server/pkg/config"
	natsservice "github.com/retawsolit/WeMeet-server/pkg/services/nats"
	"github.com/sirupsen/logrus"
)

func HandleCloseConnections() {
	if config.GetConfig() == nil {
		return
	}

	// handle to close DB connection
	if db, err := config.GetConfig().DB.DB(); err == nil {
		_ = db.Close()
	}

	// close redis
	_ = config.GetConfig().RDS.Close()

	// close nats
	natsservice.GetNatsCacheService(nil).Shutdown()
	_ = config.GetConfig().NatsConn.Drain()
	config.GetConfig().NatsConn.Close()

	// close logger
	logrus.Exit(0)
}
