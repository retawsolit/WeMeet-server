package routers

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	rr "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/retawsolit/WeMeet-server/pkg/config"
	"github.com/retawsolit/WeMeet-server/pkg/controllers"
	"github.com/retawsolit/WeMeet-server/pkg/factory"
	"github.com/retawsolit/WeMeet-server/version"
)

func New(appConfig *config.AppConfig, ctrl *factory.ApplicationControllers) *fiber.App {
	templateEngine := html.New(appConfig.Client.Path, ".html")

	if appConfig.Client.Debug {
		templateEngine.Reload(true)
		templateEngine.Debug(true)
	}

	cnf := fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Views:       templateEngine,
		AppName:     "WeMeet version: " + version.Version + " runtime: " + runtime.Version(),
	}

	if appConfig.Client.ProxyHeader != "" {
		cnf.ProxyHeader = appConfig.Client.ProxyHeader
	}

	app := fiber.New(cnf)

	if appConfig.Client.Debug {
		app.Use(logger.New())
	}
	if appConfig.Client.PrometheusConf.Enable {
		prometheus := fiberprometheus.New("WeMeet")
		prometheus.RegisterAt(app, appConfig.Client.PrometheusConf.MetricsPath)
		app.Use(prometheus.Middleware)
	}
	app.Use(rr.New())
	app.Use(cors.New(cors.Config{
		AllowMethods: "POST,GET,OPTIONS",
	}))

	clientPath := config.GetConfig().Client.Path
	// Path for wemeet-client (old meeting room UI)
	clientOldPath := "/app/client_old/dist"

	// Serve Next.js static assets (_next folder) from wemeet-ui
	app.Static("/_next", filepath.Join(clientPath, "_next"))

	// Serve static files from wemeet-ui (images, icons, etc.)
	app.Static("/apple-icon.png", filepath.Join(clientPath, "apple-icon.png"))
	app.Static("/icon.svg", filepath.Join(clientPath, "icon.svg"))
	app.Static("/icon-dark-32x32.png", filepath.Join(clientPath, "icon-dark-32x32.png"))
	app.Static("/icon-light-32x32.png", filepath.Join(clientPath, "icon-light-32x32.png"))
	app.Static("/placeholder-logo.png", filepath.Join(clientPath, "placeholder-logo.png"))
	app.Static("/placeholder-logo.svg", filepath.Join(clientPath, "placeholder-logo.svg"))
	app.Static("/placeholder-user.jpg", filepath.Join(clientPath, "placeholder-user.jpg"))
	app.Static("/placeholder.jpg", filepath.Join(clientPath, "placeholder.jpg"))
	app.Static("/placeholder.svg", filepath.Join(clientPath, "placeholder.svg"))

	// Serve wemeet-client assets (for meeting room)
	app.Static("/assets", filepath.Join(clientOldPath, "assets"))
	app.Static("/favicon.ico", filepath.Join(clientOldPath, "assets/imgs/favicon.ico"))

	// Serve root page - check for access_token to determine which UI to serve
	app.Get("/", func(c *fiber.Ctx) error {
		// If access_token is present, serve wemeet-client (meeting room UI)
		if c.Query("access_token") != "" {
			return c.SendFile(filepath.Join(clientOldPath, "index.html"))
		}
		// Otherwise serve wemeet-ui (Next.js dashboard)
		return c.SendFile(filepath.Join(clientPath, "index.html"))
	})

	// Serve login page - check if it's for wemeet-client or wemeet-ui
	app.Get("/login", func(c *fiber.Ctx) error {
		// If access_token is present, serve wemeet-client login
		if c.Query("access_token") != "" {
			loginPath := filepath.Join(clientOldPath, "login.html")
			if _, err := os.Stat(loginPath); err == nil {
				return c.SendFile(loginPath)
			}
		}
		// Otherwise serve wemeet-ui login
		loginPath := filepath.Join(clientPath, "login.html")
		if _, err := os.Stat(loginPath); err == nil {
			return c.SendFile(loginPath)
		}
		// Fallback to wemeet-ui index
		return c.SendFile(filepath.Join(clientPath, "index.html"))
	})
	app.Post("/webhook", ctrl.WebhookController.HandleWebhook)
	app.Get("/download/uploadedFile/:sid/*", ctrl.FileController.HandleDownloadUploadedFile)
	app.Get("/download/recording/:token", ctrl.RecordingController.HandleDownloadRecording)
	app.Get("/download/analytics/:token", ctrl.AnalyticsController.HandleDownloadAnalytics)
	app.Get("/healthCheck", controllers.HandleHealthCheck)

	// lti group
	lti := app.Group("/lti")
	lti.Get("/v1", ctrl.LtiV1Controller.HandleLTIV1GETREQUEST)
	lti.Post("/v1", ctrl.LtiV1Controller.HandleLTIV1Landing)
	ltiV1API := lti.Group("/v1/api", ctrl.LtiV1Controller.HandleLTIV1VerifyHeaderToken)
	ltiV1API.Post("/room/join", ctrl.LtiV1Controller.HandleLTIV1JoinRoom)
	ltiV1API.Post("/room/isActive", ctrl.LtiV1Controller.HandleLTIV1IsRoomActive)
	ltiV1API.Post("/room/end", ctrl.LtiV1Controller.HandleLTIV1EndRoom)
	ltiV1API.Post("/recording/fetch", ctrl.LtiV1Controller.HandleLTIV1FetchRecordings)
	ltiV1API.Post("/recording/download", ctrl.LtiV1Controller.HandleLTIV1GetRecordingDownloadToken)
	ltiV1API.Post("/recording/delete", ctrl.LtiV1Controller.HandleLTIV1DeleteRecordings)

	auth := app.Group("/auth", ctrl.AuthController.HandleAuthHeaderCheck)
	auth.Post("/getClientFiles", ctrl.FileController.HandleGetClientFiles)

	// for room
	room := auth.Group("/room")
	room.Post("/create", ctrl.RoomController.HandleRoomCreate)
	room.Post("/getJoinToken", ctrl.UserController.HandleGenerateJoinToken)
	room.Post("/isRoomActive", ctrl.RoomController.HandleIsRoomActive)
	room.Post("/getActiveRoomInfo", ctrl.RoomController.HandleGetActiveRoomInfo)
	room.Post("/getActiveRoomsInfo", ctrl.RoomController.HandleGetActiveRoomsInfo)
	room.Post("/endRoom", ctrl.RoomController.HandleEndRoom)
	room.Post("/fetchPastRooms", ctrl.RoomController.HandleFetchPastRooms)

	// for recording
	recording := auth.Group("/recording")
	recording.Post("/fetch", ctrl.RecordingController.HandleFetchRecordings)
	recording.Post("/recordingInfo", ctrl.RecordingController.HandleRecordingInfo)
	recording.Post("/delete", ctrl.RecordingController.HandleDeleteRecording)
	recording.Post("/getDownloadToken", ctrl.RecordingController.HandleGetDownloadToken)

	// for analytics
	analytics := auth.Group("/analytics")
	analytics.Post("/fetch", ctrl.AnalyticsController.HandleFetchAnalytics)
	analytics.Post("/delete", ctrl.AnalyticsController.HandleDeleteAnalytics)
	analytics.Post("/getDownloadToken", ctrl.AnalyticsController.HandleGetAnalyticsDownloadToken)

	// to handle different events from recorder
	recorder := auth.Group("/recorder")
	recorder.Post("/notify", ctrl.RecorderController.HandleRecorderEvents)

	// for convert BBB request to WeMeet
	bbb := app.Group("/:apiKey/bigbluebutton/api", ctrl.BBBController.HandleVerifyApiRequest)
	bbb.All("/create", ctrl.BBBController.HandleBBBCreate)
	bbb.All("/join", ctrl.BBBController.HandleBBBJoin)
	bbb.All("/isMeetingRunning", ctrl.BBBController.HandleBBBIsMeetingRunning)
	bbb.All("/getMeetingInfo", ctrl.BBBController.HandleBBBGetMeetingInfo)
	bbb.All("/getMeetings", ctrl.BBBController.HandleBBBGetMeetings)
	bbb.All("/end", ctrl.BBBController.HandleBBBEndMeetings)
	bbb.All("/getRecordings", ctrl.BBBController.HandleBBBGetRecordings)
	bbb.All("/deleteRecordings", ctrl.BBBController.HandleBBBDeleteRecordings)
	// TO-DO: in the future
	bbb.All("/updateRecordings", ctrl.BBBController.HandleBBBUpdateRecordings)
	bbb.All("/publishRecordings", ctrl.BBBController.HandleBBBPublishRecordings)

	// api group will require sending token as Authorization header value
	api := app.Group("/api", ctrl.AuthController.HandleVerifyHeaderToken)
	api.Post("/verifyToken", ctrl.AuthController.HandleVerifyToken)

	api.Post("/recording", ctrl.RecorderController.HandleRecording)
	api.Post("/rtmp", ctrl.RecorderController.HandleRTMP)
	api.Post("/endRoom", ctrl.RoomController.HandleEndRoomForAPI)
	api.Post("/changeVisibility", ctrl.RoomController.HandleChangeVisibilityForAPI)
	api.Post("/convertWhiteboardFile", ctrl.FileController.HandleConvertWhiteboardFile)
	api.Post("/externalMediaPlayer", ctrl.ExMediaController.HandleExternalMediaPlayer)
	api.Post("/externalDisplayLink", ctrl.ExDisplayController.HandleExternalDisplayLink)

	api.Post("/updateLockSettings", ctrl.UserController.HandleUpdateUserLockSetting)
	api.Post("/muteUnmuteTrack", ctrl.UserController.HandleMuteUnMuteTrack)
	api.Post("/removeParticipant", ctrl.UserController.HandleRemoveParticipant)
	api.Post("/switchPresenter", ctrl.UserController.HandleSwitchPresenter)

	// etherpad group
	etherpad := api.Group("/etherpad")
	etherpad.Post("/create", ctrl.EtherpadController.HandleCreateEtherpad)
	etherpad.Post("/cleanPad", ctrl.EtherpadController.HandleCleanPad)
	etherpad.Post("/changeStatus", ctrl.EtherpadController.HandleChangeEtherpadStatus)

	// waiting room group
	waitingRoom := api.Group("/waitingRoom")
	waitingRoom.Post("/approveUsers", ctrl.WaitingRoomController.HandleApproveUsers)
	waitingRoom.Post("/updateMsg", ctrl.WaitingRoomController.HandleUpdateWaitingRoomMessage)

	// polls group
	polls := api.Group("/polls")
	polls.Post("/activate", ctrl.PollsController.HandleActivatePolls)
	polls.Post("/create", ctrl.PollsController.HandleCreatePoll)
	polls.Get("/listPolls", ctrl.PollsController.HandleListPolls)
	polls.Get("/pollsStats", ctrl.PollsController.HandleGetPollsStats)
	polls.Get("/countTotalResponses/:pollId", ctrl.PollsController.HandleCountPollTotalResponses)
	polls.Get("/userSelectedOption/:pollId/:userId", ctrl.PollsController.HandleUserSelectedOption)
	polls.Get("/pollResponsesDetails/:pollId", ctrl.PollsController.HandleGetPollResponsesDetails)
	polls.Get("/pollResponsesResult/:pollId", ctrl.PollsController.HandleGetResponsesResult)
	polls.Post("/submitResponse", ctrl.PollsController.HandleUserSubmitResponse)
	polls.Post("/closePoll", ctrl.PollsController.HandleClosePoll)

	// breakout room group
	breakoutRoom := api.Group("/breakoutRoom")
	breakoutRoom.Post("/create", ctrl.BreakoutRoomController.HandleCreateBreakoutRooms)
	breakoutRoom.Post("/join", ctrl.BreakoutRoomController.HandleJoinBreakoutRoom)
	breakoutRoom.Get("/listRooms", ctrl.BreakoutRoomController.HandleGetBreakoutRooms)
	breakoutRoom.Get("/myRooms", ctrl.BreakoutRoomController.HandleGetMyBreakoutRooms)
	breakoutRoom.Post("/increaseDuration", ctrl.BreakoutRoomController.HandleIncreaseBreakoutRoomDuration)
	breakoutRoom.Post("/sendMsg", ctrl.BreakoutRoomController.HandleSendBreakoutRoomMsg)
	breakoutRoom.Post("/endRoom", ctrl.BreakoutRoomController.HandleEndBreakoutRoom)
	breakoutRoom.Post("/endAllRooms", ctrl.BreakoutRoomController.HandleEndBreakoutRooms)

	// Ingress
	ingress := api.Group("/ingress")
	ingress.Post("/create", ctrl.IngressController.HandleCreateIngress)

	// Speech services
	speech := api.Group("/speechServices")
	speech.Post("/serviceStatus", ctrl.SpeechToTextController.HandleSpeechToTextTranslationServiceStatus)
	speech.Post("/azureToken", ctrl.SpeechToTextController.HandleGenerateAzureToken)
	speech.Post("/userStatus", ctrl.SpeechToTextController.HandleSpeechServiceUserStatus)
	speech.Post("/renewToken", ctrl.SpeechToTextController.HandleRenewAzureToken)

	// for resumable.js need both GET and POST  methods.
	// https://github.com/23/resumable.js#how-do-i-set-it-up-with-my-server
	api.Get("/fileUpload", ctrl.FileController.HandleFileUpload)
	api.Post("/fileUpload", ctrl.FileController.HandleFileUpload)
	// as resumable.js will upload multiple parts of the file in different request
	// merging request should be sent from another request
	// otherwise hard to do it concurrently
	api.Post("/uploadedFileMerge", ctrl.FileController.HandleUploadedFileMerge)
	api.Post("/uploadBase64EncodedData", ctrl.FileController.HandleUploadBase64EncodedData)

	// SPA fallback: serve HTML files for UI routes, but not for API routes
	app.Use(func(c *fiber.Ctx) error {
		path := c.Path()

		// Skip API routes, auth routes, download routes, webhook, healthCheck, BBB, and LTI
		if strings.HasPrefix(path, "/api/") ||
			strings.HasPrefix(path, "/auth/") ||
			strings.HasPrefix(path, "/room/") ||
			strings.HasPrefix(path, "/lti/") ||
			strings.HasPrefix(path, "/webhook") ||
			strings.HasPrefix(path, "/download/") ||
			strings.HasPrefix(path, "/healthCheck") ||
			strings.HasPrefix(path, "/:apiKey/") ||
			strings.HasPrefix(path, "/_next/") ||
			strings.HasPrefix(path, "/assets/") {
			return c.Status(fiber.StatusNotFound).SendString("not found")
		}

		// Try to serve the corresponding HTML file for the route
		// e.g., /dashboard -> dashboard.html, /login -> login.html
		routePath := strings.TrimPrefix(path, "/")
		if routePath == "" {
			routePath = "index"
		}

		htmlPath := filepath.Join(clientPath, routePath+".html")
		// Check if file exists before trying to serve it
		if _, err := os.Stat(htmlPath); err == nil {
			// File exists, serve it
			return c.SendFile(htmlPath)
		}

		// Fallback to index.html for SPA routing (for dynamic routes)
		return c.SendFile(filepath.Join(clientPath, "index.html"))
	})

	return app
}
