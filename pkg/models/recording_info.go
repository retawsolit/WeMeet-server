package models

import (
	"errors"

	"github.com/retawsolit/wemeet-protocol/wemeet"
)

func (m *RecordingModel) FetchRecordings(r *wemeet.FetchRecordingsReq) (*wemeet.FetchRecordingsResult, error) {
	if r.Limit <= 0 {
		r.Limit = 20
	}
	if r.OrderBy == "" {
		r.OrderBy = "DESC"
	}

	data, total, err := m.ds.GetRecordings(r.RoomIds, uint64(r.From), uint64(r.Limit), &r.OrderBy)
	if err != nil {
		return nil, err
	}
	var recordings []*wemeet.RecordingInfo
	for _, v := range data {
		recording := &wemeet.RecordingInfo{
			RecordId:         v.RecordID,
			RoomId:           v.RoomID,
			RoomSid:          v.RoomSid.String,
			FilePath:         v.FilePath,
			FileSize:         float32(v.Size),
			CreationTime:     v.CreationTime,
			RoomCreationTime: v.RoomCreationTime,
		}
		recordings = append(recordings, recording)
	}

	result := &wemeet.FetchRecordingsResult{
		TotalRecordings: total,
		From:            r.From,
		Limit:           r.Limit,
		OrderBy:         r.OrderBy,
		RecordingsList:  recordings,
	}

	return result, nil
}

// FetchRecording to get single recording information from DB
func (m *RecordingModel) FetchRecording(recordId string) (*wemeet.RecordingInfo, error) {
	v, err := m.ds.GetRecording(recordId)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, errors.New("no info found")
	}
	recording := &wemeet.RecordingInfo{
		RecordId:         v.RecordID,
		RoomId:           v.RoomID,
		RoomSid:          v.RoomSid.String,
		FilePath:         v.FilePath,
		FileSize:         float32(v.Size),
		CreationTime:     v.CreationTime,
		RoomCreationTime: v.RoomCreationTime,
	}

	return recording, nil
}

func (m *RecordingModel) RecordingInfo(req *wemeet.RecordingInfoReq) (*wemeet.RecordingInfoRes, error) {
	recording, err := m.FetchRecording(req.RecordId)
	if err != nil {
		return nil, err
	}

	pastRoomInfo := new(wemeet.PastRoomInfo)
	// SID can't be null, so we'll check before
	if recording.GetRoomSid() != "" {
		if roomInfo, err := m.ds.GetRoomInfoBySid(recording.GetRoomSid(), nil); err == nil && roomInfo != nil {
			pastRoomInfo = &wemeet.PastRoomInfo{
				RoomTitle:          roomInfo.RoomTitle,
				RoomId:             roomInfo.RoomId,
				RoomSid:            roomInfo.Sid,
				JoinedParticipants: roomInfo.JoinedParticipants,
				WebhookUrl:         roomInfo.WebhookUrl,
				Created:            roomInfo.Created.Format("2006-01-02 15:04:05"),
				Ended:              roomInfo.Ended.Format("2006-01-02 15:04:05"),
			}
			if an, err := m.ds.GetAnalyticByRoomTableId(roomInfo.ID); err == nil && an != nil {
				pastRoomInfo.AnalyticsFileId = an.FileID
			}
		}
	}

	return &wemeet.RecordingInfoRes{
		Status:        true,
		Msg:           "success",
		RecordingInfo: recording,
		RoomInfo:      pastRoomInfo,
	}, nil
}
