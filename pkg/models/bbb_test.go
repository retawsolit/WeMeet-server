package models

import (
	"testing"

	"github.com/retawsolit/wemeet-protocol/bbbapiwrapper"
)

func TestBBBApiWrapperModel_GetRecordings(t *testing.T) {
	bbbm := NewBBBApiWrapperModel(nil, nil, nil)
	recordings, pag, err := bbbm.GetRecordings("https://demo.wemeet.com", &bbbapiwrapper.GetRecordingsReq{
		MeetingID: roomId,
	})
	if err != nil {
		t.Error(err)
	}
	if len(recordings) == 0 {
		t.Error("should contains some data but got empty")
	}

	t.Logf("%+v, %+v", recordings[0], *pag)

	recordings, pag, err = bbbm.GetRecordings("https://demo.wemeet.com", &bbbapiwrapper.GetRecordingsReq{
		RecordID: recordId,
	})
	if err != nil {
		t.Error(err)
	}
	if len(recordings) == 0 {
		t.Error("should contains some data but got empty")
	}

	t.Logf("%+v, %+v", recordings[0], *pag)
}
