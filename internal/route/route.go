package route

import (
	"net/http"

	"github.com/go-kit/log"

	"github.com/gin-gonic/gin"
	"gitlab.wlink.com.np/nettv-webhook/internal/controller"
	"gitlab.wlink.com.np/nettv-webhook/messagebroker"
)

func SetupRoutes(msgBrk messagebroker.MessageBrokerService, logger log.Logger) *gin.Engine {
	httpRouter := gin.Default()

	httpRouter.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "welcome to simple webhook")
	})

	webHookController := controller.NewWebHookController(logger, msgBrk)

	httpRouter.POST("/webhook", webHookController.HandleWebhook)

	return httpRouter

}
