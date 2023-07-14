package http

import (
	"github.com/gin-gonic/gin"
)

type ServerURI string

const ServerURISlack ServerURI = "/slack"
const ServerURISms ServerURI = "/sms"
const ServerURIEmail ServerURI = "/email"

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	router.POST(string(ServerURISlack), handler.HandleSlackNotification)
	router.POST(string(ServerURISms), handler.HandleSMSNotification)
	router.POST(string(ServerURIEmail), handler.HandleEmailNotification)
}
