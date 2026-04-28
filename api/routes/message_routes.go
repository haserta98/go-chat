package routes

import (
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal/handler"
	"github.com/haserta98/go-rest/internal/middleware"
	"github.com/haserta98/go-rest/internal/service"
)

func initMessageEndpoints(appCtx *cmd.AppContext) {
	httpServer := appCtx.GetHTTPServer()
	
	// MessageService is already registered in appCtx
	msgSvc := appCtx.GetService("MessageService").(*service.MessageService)
	messageHandler := handler.NewMessageHandler(msgSvc)

	authMiddleware := middleware.NewAuthMiddleware(appCtx.GetRedisClient())
	protected := httpServer.GetInstance().Group("/messages", authMiddleware)
	protected.Get("/:otherUserID", messageHandler.GetMessagesBetween)
	protected.Get("/group/:groupID", messageHandler.GetGroupMessages)
}
