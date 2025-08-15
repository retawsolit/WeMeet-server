package models

import (
	"time"
)

type RoomInfo struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	RoomTitle          string    `gorm:"size:255;not null" json:"room_title"`
	RoomId             string    `gorm:"size:64;not null;uniqueIndex" json:"room_id"`
	Sid                string    `gorm:"size:64;not null;unique" json:"sid"`
	JoinedParticipants int       `gorm:"default:0" json:"joined_participants"`
	IsRunning          bool      `gorm:"default:false" json:"is_running"`
	IsRecording        bool      `gorm:"default:false" json:"is_recording"`
	RecorderID         string    `gorm:"size:36;default:''" json:"recorder_id"`
	IsActiveRtmp       bool      `gorm:"default:false" json:"is_active_rtmp"`
	RtmpNodeID         string    `gorm:"size:36;default:''" json:"rtmp_node_id"`
	WebhookUrl         string    `gorm:"size:255;default:''" json:"webhook_url"`
	IsBreakoutRoom     bool      `gorm:"default:false" json:"is_breakout_room"`
	ParentRoomID       string    `gorm:"size:64;default:''" json:"parent_room_id"`
	CreationTime       int64     `gorm:"default:0" json:"creation_time"`
	Created            time.Time `gorm:"autoCreateTime" json:"created"`
	Ended              time.Time `json:"ended"`
	Modified           time.Time `gorm:"autoUpdateTime" json:"modified"`
}

func (RoomInfo) TableName() string {
	return "pnm_room_info"
}
