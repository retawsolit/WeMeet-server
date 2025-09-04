package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/retawsolit/WeMeet-server/pkg/models"
	"github.com/retawsolit/wemeet-protocol/utils"
	"github.com/retawsolit/wemeet-protocol/wemeet"
)

// RecordingController holds dependencies for recording-related handlers.
type RecordingController struct {
	RecordingModel *models.RecordingModel
}

// NewRecordingController creates a new RecordingController.
func NewRecordingController(m *models.RecordingModel) *RecordingController {
	return &RecordingController{
		RecordingModel: m,
	}
}

// HandleFetchRecordings handles fetching recordings.
func (rc *RecordingController) HandleFetchRecordings(c *fiber.Ctx) error {
	req := new(wemeet.FetchRecordingsReq)
	if err := parseAndValidateRequest(c.Body(), req); err != nil {
		return utils.SendCommonProtoJsonResponse(c, false, err.Error())
	}

	result, err := rc.RecordingModel.FetchRecordings(req)
	if err != nil {
		return utils.SendCommonProtoJsonResponse(c, false, err.Error())
	}
	if result.GetTotalRecordings() == 0 {
		return utils.SendCommonProtoJsonResponse(c, false, "no recordings found")
	}

	r := &wemeet.FetchRecordingsRes{
		Status: true,
		Msg:    "success",
		Result: result,
	}
	return utils.SendProtoJsonResponse(c, r)
}

// HandleRecordingInfo handles fetching information for a single recording.
func (rc *RecordingController) HandleRecordingInfo(c *fiber.Ctx) error {
	req := new(wemeet.RecordingInfoReq)
	if err := parseAndValidateRequest(c.Body(), req); err != nil {
		return utils.SendCommonProtoJsonResponse(c, false, err.Error())
	}

	result, err := rc.RecordingModel.RecordingInfo(req)
	if err != nil {
		return utils.SendCommonProtoJsonResponse(c, false, err.Error())
	}

	return utils.SendProtoJsonResponse(c, result)
}

// HandleDeleteRecording handles deleting a recording.
func (rc *RecordingController) HandleDeleteRecording(c *fiber.Ctx) error {
	req := new(wemeet.DeleteRecordingReq)
	if err := parseAndValidateRequest(c.Body(), req); err != nil {
		return utils.SendCommonProtoJsonResponse(c, false, err.Error())
	}

	err := rc.RecordingModel.DeleteRecording(req)
	if err != nil {
		return utils.SendCommonProtoJsonResponse(c, false, err.Error())
	}

	return utils.SendCommonProtoJsonResponse(c, true, "success")
}

// HandleGetDownloadToken handles generating a download token for a recording.
func (rc *RecordingController) HandleGetDownloadToken(c *fiber.Ctx) error {
	req := new(wemeet.GetDownloadTokenReq)
	if err := parseAndValidateRequest(c.Body(), req); err != nil {
		return utils.SendCommonProtoJsonResponse(c, false, err.Error())
	}

	token, err := rc.RecordingModel.GetDownloadToken(req)
	if err != nil {
		return utils.SendCommonProtoJsonResponse(c, false, err.Error())
	}

	r := &wemeet.GetDownloadTokenRes{
		Status: true,
		Msg:    "success",
		Token:  &token,
	}
	return utils.SendProtoJsonResponse(c, r)
}

// HandleDownloadRecording handles downloading a recording file.
func (rc *RecordingController) HandleDownloadRecording(c *fiber.Ctx) error {
	token := c.Params("token")

	if len(token) == 0 {
		return c.Status(fiber.StatusUnauthorized).SendString("token require or invalid url")
	}

	file, status, err := rc.RecordingModel.VerifyRecordingToken(token)
	if err != nil {
		return c.Status(status).SendString(err.Error())
	}

	c.Attachment(file)
	return c.SendFile(file, false)
}
