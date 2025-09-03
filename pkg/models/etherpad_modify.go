package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/retawsolit/!we!meet-protocol/wemeet backup moi"
	log "github.com/sirupsen/logrus"
)

func (m *EtherpadModel) ChangeEtherpadStatus(r *wemeet.ChangeEtherpadStatusReq) error {
	meta, err := m.natsService.GetRoomMetadataStruct(r.RoomId)
	if err != nil {
		return err
	}
	if meta == nil {
		return errors.New("invalid nil room metadata information")
	}

	meta.RoomFeatures.SharedNotePadFeatures.IsActive = r.IsActive
	err = m.natsService.UpdateAndBroadcastRoomMetadata(r.RoomId, meta)
	if err != nil {
		log.Errorln(err)
	}

	// send analytics
	val := wemeet.AnalyticsStatus_ANALYTICS_STATUS_STARTED.String()
	d := &wemeet.AnalyticsDataMsg{
		EventType: wemeet.AnalyticsEventType_ANALYTICS_EVENT_TYPE_ROOM,
		EventName: wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_ETHERPAD_STATUS,
		RoomId:    r.RoomId,
		HsetValue: &val,
	}
	if !r.IsActive {
		val = wemeet.AnalyticsStatus_ANALYTICS_STATUS_ENDED.String()
		d.EventName = wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_ETHERPAD_STATUS
		d.HsetValue = &val
	}
	m.analyticsModel.HandleEvent(d)

	return err
}

func (m *EtherpadModel) addPadToRoomMetadata(roomId string, c *wemeet.CreateEtherpadSessionRes) error {
	meta, err := m.natsService.GetRoomMetadataStruct(roomId)
	if err != nil {
		return err
	}
	if meta == nil {
		return errors.New("invalid room information")
	}

	f := &wemeet.SharedNotePadFeatures{
		AllowedSharedNotePad: meta.RoomFeatures.SharedNotePadFeatures.AllowedSharedNotePad,
		IsActive:             true,
		NodeId:               m.NodeId,
		Host:                 m.Host,
		NotePadId:            *c.PadId,
		ReadOnlyPadId:        *c.ReadonlyPadId,
	}
	meta.RoomFeatures.SharedNotePadFeatures = f

	err = m.natsService.UpdateAndBroadcastRoomMetadata(roomId, meta)
	if err != nil {
		log.Errorln(err)
	}

	// send analytics
	val := wemeet.AnalyticsStatus_ANALYTICS_STATUS_STARTED.String()
	m.analyticsModel.HandleEvent(&wemeet.AnalyticsDataMsg{
		EventType: wemeet.AnalyticsEventType_ANALYTICS_EVENT_TYPE_ROOM,
		EventName: wemeet.AnalyticsEvents_ANALYTICS_EVENT_ROOM_ETHERPAD_STATUS,
		RoomId:    roomId,
		HsetValue: &val,
	})

	return err
}

func (m *EtherpadModel) postToEtherpad(method string, vals url.Values) (*EtherpadHttpRes, error) {
	if m.NodeId == "" {
		return nil, errors.New("no notepad nodeId found")
	}
	token, err := m.getAccessToken()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	en := vals.Encode()
	endPoint := fmt.Sprintf("%s/api/%s/%s?%s", m.Host, APIVersion, method, en)

	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("error code: " + res.Status)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	mar := new(EtherpadHttpRes)
	err = json.Unmarshal(body, mar)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}

	return mar, nil
}

func (m *EtherpadModel) getAccessToken() (string, error) {
	token, _ := m.natsService.GetEtherpadToken(m.NodeId)
	if token != "" {
		return token, nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", m.ClientId)
	data.Set("client_secret", m.ClientSecret)
	encodedData := data.Encode()

	client := &http.Client{}
	urlPath := fmt.Sprintf("%s/oidc/token", m.Host)

	req, err := http.NewRequest("POST", urlPath, strings.NewReader(encodedData))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	vals := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = json.Unmarshal(body, &vals)
	if err != nil {
		return "", err
	}

	if vals.AccessToken == "" {
		return "", errors.New("can not get access_token value")
	}

	// we'll store the value with expiry of 30-minute max
	err = m.natsService.AddEtherpadToken(m.NodeId, vals.AccessToken, time.Minute*30)
	if err != nil {
		log.Errorln(err)
	}

	return vals.AccessToken, nil
}
