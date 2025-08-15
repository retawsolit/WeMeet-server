package models

import (
	"time"
)

type Recording struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	RecordID         string    `gorm:"size:64;not null;unique" json:"record_id"`
	RoomID           string    `gorm:"size:64;not null" json:"room_id"`
	RoomSid          *string   `gorm:"size:64" json:"room_sid"`
	RecorderID       string    `gorm:"size:36;not null" json:"recorder_id"`
	FilePath         string    `gorm:"size:255;not null" json:"file_path"`
	Size             float64   `gorm:"not null" json:"size"`
	Published        bool      `gorm:"default:true" json:"published"`
	CreationTime     int64     `gorm:"default:0" json:"creation_time"`
	RoomCreationTime int64     `gorm:"default:0" json:"room_creation_time"`
	Created          time.Time `gorm:"autoCreateTime" json:"created"`
	Modified         time.Time `gorm:"autoUpdateTime" json:"modified"`

	// Foreign key relationship
	Room *RoomInfo `gorm:"foreignKey:RoomSid;references:Sid" json:"room,omitempty"`
}

func (Recording) TableName() string {
	return "pnm_recordings"
}
