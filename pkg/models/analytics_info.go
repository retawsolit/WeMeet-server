package models

import (
	"errors"
)

func (m *AnalyticsModel) FetchAnalytics(r *wemeet.FetchAnalyticsReq) (*wemeet.FetchAnalyticsResult, error) {
	if r.Limit <= 0 {
		r.Limit = 20
	}
	if r.OrderBy == "" {
		r.OrderBy = "DESC"
	}
	data, total, err := m.ds.GetAnalytics(r.RoomIds, uint64(r.From), uint64(r.Limit), &r.OrderBy)
	if err != nil {
		return nil, err
	}

	var analytics []*wemeet.AnalyticsInfo
	for _, v := range data {
		analytic := &wemeet.AnalyticsInfo{
			RoomId:           v.RoomID,
			FileId:           v.FileID,
			FileSize:         v.FileSize,
			FileName:         v.FileName,
			CreationTime:     v.CreationTime,
			RoomCreationTime: v.RoomCreationTime,
		}
		analytics = append(analytics, analytic)
	}

	result := &wemeet.FetchAnalyticsResult{
		TotalAnalytics: total,
		From:           r.From,
		Limit:          r.Limit,
		OrderBy:        r.OrderBy,
		AnalyticsList:  analytics,
	}

	return result, nil
}

func (m *AnalyticsModel) fetchAnalytic(fileId string) (*wemeet.AnalyticsInfo, error) {
	v, err := m.ds.GetAnalyticByFileId(fileId)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, errors.New("no info found")
	}
	analytic := &wemeet.AnalyticsInfo{
		RoomId:           v.RoomID,
		FileId:           v.FileID,
		FileSize:         v.FileSize,
		FileName:         v.FileName,
		CreationTime:     v.CreationTime,
		RoomCreationTime: v.RoomCreationTime,
	}

	return analytic, nil
}

func (m *AnalyticsModel) getAnalyticByRoomTableId(roomTableId uint64) (*wemeet.AnalyticsInfo, error) {
	v, err := m.ds.GetAnalyticByRoomTableId(roomTableId)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, errors.New("no info found")
	}
	analytic := &wemeet.AnalyticsInfo{
		RoomId:           v.RoomID,
		FileId:           v.FileID,
		FileName:         v.FileName,
		FileSize:         v.FileSize,
		CreationTime:     v.CreationTime,
		RoomCreationTime: v.RoomCreationTime,
	}

	return analytic, nil
}
