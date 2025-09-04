package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/retawsolit/wemeet-protocol/wemeet"
)

var recordingModel *RecordingModel

func TestAuthRecording_FetchRecordings(t *testing.T) {
	recordingModel = NewRecordingModel(nil, nil, nil)

	result, err := recordingModel.FetchRecordings(&wemeet.FetchRecordingsReq{
		RoomIds: []string{roomId},
	})
	if err != nil {
		t.Error(err)
	}

	if len(result.RecordingsList) == 0 {
		t.Error("should contains some data but got empty")
	}

	t.Logf("%+v", result)
}

func TestAnalyticsAuthModel_FetchRecording(t *testing.T) {
	result, err := recordingModel.FetchRecording(recordId)
	if err != nil {
		t.Error(err)
	}

	if result == nil {
		t.Error("should contains some data but got empty")
	}

	t.Logf("%+v", result)

	_, err = recordingModel.FetchRecording(fmt.Sprintf("%d", time.Now().UnixMilli()))
	if err == nil {
		t.Error("should got not found error")
	}
}

func TestAnalyticsAuthModel_RecordingInfo(t *testing.T) {
	result, err := recordingModel.RecordingInfo(&wemeet.RecordingInfoReq{
		RecordId: recordId,
	})
	if err != nil {
		t.Error(err)
	}
	if result == nil {
		t.Error("should contains some data but got empty")
		return
	}

	if result.RoomInfo.RoomSid != sid {
		t.Errorf("sid of %s should match with our %s", result.RoomInfo.RoomSid, sid)
	}

	t.Logf("%+v", result)
}

func TestAnalyticsAuthModel_DeleteRecording(t *testing.T) {
	err := recordingModel.DeleteRecording(&wemeet.DeleteRecordingReq{
		RecordId: recordId,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestAnalyticsAuthModel_CreateAndVerifyToken(t *testing.T) {
	token, err := recordingModel.CreateTokenForDownload("test.mp4")
	if err != nil {
		t.Error(err)
	}
	_, res, err := recordingModel.VerifyRecordingToken(token)
	if err == nil {
		t.Error("should not found the file")
		return
	}

	if res != fiber.StatusNotFound {
		t.Errorf("should get response: %d", fiber.StatusNotFound)
	}

	t.Logf("%+v, response: %d", err, fiber.StatusNotFound)
}
