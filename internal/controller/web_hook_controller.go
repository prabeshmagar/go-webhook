package controller

import (
	"io"

	"github.com/go-kit/log"

	"github.com/gin-gonic/gin"
	"gitlab.wlink.com.np/nettv-webhook/messagebroker"
)

type WebHookController interface {
	HandleWebhook(c *gin.Context)
}

type webhookController struct {
	logger log.Logger
	msgBrk messagebroker.MessageBrokerService
}

func NewWebHookController(logger log.Logger, mb messagebroker.MessageBrokerService) *webhookController {
	return &webhookController{
		logger: logger,
		msgBrk: mb,
	}
}

func (w webhookController) HandleWebhook(c *gin.Context) {
	data, _ := io.ReadAll(c.Request.Body)
	w.msgBrk.Publish(data)
}
