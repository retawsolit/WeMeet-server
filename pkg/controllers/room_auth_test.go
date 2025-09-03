package controllers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/retawsolit/WeMeet-protocol/wemeet"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

func setupApp() *fiber.App {
	app := fiber.New()
	room := app.Group("/room")
	room.Post("/create", HandleRoomCreate)
	room.Post("/getJoinToken", HandleGenerateJoinToken)
	room.Post("/isRoomActive", HandleIsRoomActive)
	room.Post("/getActiveRoomInfo", HandleGetActiveRoomInfo)
	room.Post("/getActiveRoomsInfo", HandleGetActiveRoomsInfo)
	room.Post("/endRoom", HandleEndRoom)
	room.Post("/fetchPastRooms", HandleFetchPastRooms)

	// some APIs
	api := app.Group("/api", HandleVerifyHeaderToken)
	api.Post("/verifyToken", HandleVerifyToken)

	// others
	app.Post("/webhook", HandleWebhook)
	return app
}

func TestHandleRoomCreate(t *testing.T) {
	app := setupApp()
	reqBody := &wemeet.CreateRoomReq{
		RoomId: roomId,
		Metadata: &wemeet.RoomMetadata{
			RoomTitle:    "Test Room",
			RoomFeatures: &wemeet.RoomCreateFeatures{},
		},
	}

	bodyBytes, err := protojson.Marshal(reqBody)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/room/create", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Read and unmarshal response
	respBody := new(wemeet.CreateRoomRes)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	assert.NoError(t, err)

	err = protojson.Unmarshal(buf.Bytes(), respBody)
	assert.NoError(t, err)

	// Compare expected values
	assert.True(t, respBody.Status)
	assert.Equal(t, "success", respBody.Msg)
	assert.NotNil(t, respBody.RoomInfo)
	assert.Equal(t, roomId, respBody.RoomInfo.RoomId)
}

func TestHandleGetJoinToken(t *testing.T) {
	app := setupApp()
	reqBody := &wemeet.GenerateTokenReq{
		RoomId: roomId,
		UserInfo: &wemeet.UserInfo{
			UserId: userId,
			Name:   "Test User",
			UserMetadata: &wemeet.UserMetadata{
				IsAdmin: true,
			},
		},
	}

	bodyBytes, err := protojson.Marshal(reqBody)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/room/getJoinToken", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Read and unmarshal response
	respBody := new(wemeet.GenerateTokenRes)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	assert.NoError(t, err)

	err = protojson.Unmarshal(buf.Bytes(), respBody)
	assert.NoError(t, err)

	// Compare expected values
	assert.True(t, respBody.Status)
	assert.Equal(t, "success", respBody.Msg)
	assert.NotNil(t, respBody.Token)

	if respBody.Token != nil {
		testValidateJoinToken(t, *respBody.Token)
	}
}

func TestHandleIsRoomActive(t *testing.T) {
	app := setupApp()
	reqBody := &wemeet.IsRoomActiveReq{
		RoomId: roomId,
	}

	bodyBytes, err := protojson.Marshal(reqBody)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/room/isRoomActive", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Read and unmarshal response
	respBody := new(wemeet.IsRoomActiveRes)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	assert.NoError(t, err)

	err = protojson.Unmarshal(buf.Bytes(), respBody)
	assert.NoError(t, err)
	// Compare expected values
	assert.True(t, respBody.Status)
	assert.Equal(t, "room is active", respBody.Msg)
	assert.True(t, respBody.IsActive)
}

func TestHandleGetActiveRoomInfo(t *testing.T) {
	app := setupApp()
	reqBody := &wemeet.GetActiveRoomInfoReq{
		RoomId: roomId,
	}

	bodyBytes, err := protojson.Marshal(reqBody)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/room/getActiveRoomInfo", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Read and unmarshal response
	respBody := new(wemeet.GetActiveRoomInfoRes)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	assert.NoError(t, err)

	err = protojson.Unmarshal(buf.Bytes(), respBody)
	assert.NoError(t, err)
	// Compare expected values
	assert.True(t, respBody.Status)
	assert.Equal(t, "success", respBody.Msg)
	assert.NotNil(t, respBody.Room)
	assert.NotNil(t, respBody.Room.RoomInfo)
	assert.Equal(t, roomId, respBody.Room.RoomInfo.RoomId)
}

func TestHandleGetActiveRoomsInfo(t *testing.T) {
	app := setupApp()
	req := httptest.NewRequest("POST", "/room/getActiveRoomsInfo", nil)

	// Send request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Read and unmarshal response
	respBody := new(wemeet.GetActiveRoomsInfoRes)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	assert.NoError(t, err)

	err = protojson.Unmarshal(buf.Bytes(), respBody)
	assert.NoError(t, err)
	// Compare expected values
	assert.True(t, respBody.Status)
	assert.Equal(t, "success", respBody.Msg)
	assert.GreaterOrEqual(t, len(respBody.Rooms), 1)
}

func TestHandleEndRoom(t *testing.T) {
	app := setupApp()
	reqBody := &wemeet.RoomEndReq{
		RoomId: roomId,
	}

	bodyBytes, err := protojson.Marshal(reqBody)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/room/endRoom", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Read and unmarshal response
	respBody := new(wemeet.RoomEndRes)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	assert.NoError(t, err)

	err = protojson.Unmarshal(buf.Bytes(), respBody)
	assert.NoError(t, err)
	// Compare expected values
	assert.True(t, respBody.Status)
	assert.Equal(t, "success", respBody.Msg)
}

func TestHandleFetchPastRooms(t *testing.T) {
	app := setupApp()
	reqBody := &wemeet.FetchPastRoomsReq{
		RoomIds: []string{roomId},
	}

	bodyBytes, err := protojson.Marshal(reqBody)
	assert.NoError(t, err)
	req := httptest.NewRequest("POST", "/room/fetchPastRooms", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Read and unmarshal response
	respBody := new(wemeet.FetchPastRoomsRes)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	assert.NoError(t, err)

	err = protojson.Unmarshal(buf.Bytes(), respBody)
	assert.NoError(t, err)
	// Compare expected values
	assert.True(t, respBody.Status)
	assert.Equal(t, "success", respBody.Msg)
	assert.NotNil(t, respBody.Result)
	assert.NotNil(t, respBody.Result.RoomsList)
	assert.GreaterOrEqual(t, len(respBody.Result.RoomsList), 1)
}
