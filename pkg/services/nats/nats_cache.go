package natsservice

import (
	"context"
	"strconv"
	"sync"

	"github.com/retawsolit/WeMeet-server/pkg/config"
	"github.com/retawsolit/wemeet-protocol/wemeet"
	log "github.com/sirupsen/logrus"
)

var (
	defaultNatsCacheService *NatsCacheService
	initCacheOnce           sync.Once
)

type CachedRoomEntry struct {
	RoomInfo *wemeet.NatsKvRoomInfo
}

type CachedRoomUserStatusEntry struct {
	Status   string
	Revision uint64
}

type CachedUserInfoEntry struct {
	UserInfo   *wemeet.NatsKvUserInfo
	LastPingAt uint64
}

type NatsCacheService struct {
	// Global context for all long-lived watchers in this service
	serviceCtx    context.Context
	serviceCancel context.CancelFunc

	roomLock       sync.RWMutex
	roomsInfoStore map[string]CachedRoomEntry

	userLock             sync.RWMutex
	roomUsersStatusStore map[string]map[string]CachedRoomUserStatusEntry
	roomUsersInfoStore   map[string]map[string]CachedUserInfoEntry
}

func InitNatsCacheService(app *config.AppConfig) {
	initCacheOnce.Do(func() {
		if app == nil {
			app = config.GetConfig()
		}
		if app.JetStream == nil {
			log.Fatal("NATS JetStream not provided to InitNatsCacheService")
		}

		ctx, cancel := context.WithCancel(context.Background())
		defaultNatsCacheService = &NatsCacheService{
			serviceCtx:           ctx,
			serviceCancel:        cancel,
			roomsInfoStore:       make(map[string]CachedRoomEntry),
			roomUsersStatusStore: make(map[string]map[string]CachedRoomUserStatusEntry),
			roomUsersInfoStore:   make(map[string]map[string]CachedUserInfoEntry),
		}
	})
}

// GetNatsCacheService returns the singleton instance.
func GetNatsCacheService(app *config.AppConfig) *NatsCacheService {
	if defaultNatsCacheService == nil {
		InitNatsCacheService(app)
	}
	return defaultNatsCacheService
}

// Shutdown gracefully stops all watchers.
func (ncs *NatsCacheService) Shutdown() {
	log.Info("Shutting down NATS Cache Service...")
	ncs.serviceCancel() // Signals all watchers started with ncs.serviceCtx to stop
	log.Info("NATS Cache Service shutdown complete.")
}

func (ncs *NatsCacheService) convertTextToUint64(text string) uint64 {
	value, _ := strconv.ParseUint(text, 10, 64)
	return value
}
