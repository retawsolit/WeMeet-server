package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/retawsolit/WeMeet-protocol/utils"
	"github.com/retawsolit/WeMeet-protocol/wemeet"
	"github.com/retawsolit/WeMeet-server/pkg/models"
	"google.golang.org/protobuf/proto"
)

// ExMediaController holds dependencies for external media player handlers.
type ExMediaController struct {
	ExMediaModel *models.ExMediaModel
}

// NewExMediaController creates a new ExMediaController.
func NewExMediaController(emm *models.ExMediaModel) *ExMediaController {
	return &ExMediaController{
		ExMediaModel: emm,
	}
}

// HandleExternalMediaPlayer handles external media player actions.
func (emc *ExMediaController) HandleExternalMediaPlayer(c *fiber.Ctx) error {
	isAdmin := c.Locals("isAdmin")
	roomId := c.Locals("roomId")
	requestedUserId := c.Locals("requestedUserId")

	if !isAdmin.(bool) {
		return utils.SendCommonProtobufResponse(c, false, "only admin can perform this task")
	}

	rid := roomId.(string)
	if rid == "" {
		return utils.SendCommonProtobufResponse(c, false, "roomId required")
	}

	req := new(wemeet.ExternalMediaPlayerReq)
	err := proto.Unmarshal(c.Body(), req)
	if err != nil {
		return utils.SendCommonProtobufResponse(c, false, err.Error())
	}

	req.RoomId = rid
	req.UserId = requestedUserId.(string)
	err = emc.ExMediaModel.HandleTask(req)
	if err != nil {
		return utils.SendCommonProtobufResponse(c, false, err.Error())
	}

	return utils.SendCommonProtobufResponse(c, true, "success")
}
