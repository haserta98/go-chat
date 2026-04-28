package routes

import (
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal/handler"
	"github.com/haserta98/go-rest/internal/middleware"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/haserta98/go-rest/internal/service"
	"github.com/haserta98/go-rest/internal/ws"
)

func initGroupEndpoints(appCtx *cmd.AppContext) {
	httpServer := appCtx.GetHTTPServer()
	wsGateway := appCtx.GetService("WsGateway").(*ws.WsGateway)
	groupService := service.NewGroupService(appCtx.GetRepository("Group").(*repository.GroupRepository), wsGateway)
	groupService.Start()
	groupHandler := handler.NewGroupHandler(groupService)

	authMiddleware := middleware.NewAuthMiddleware(appCtx.GetRedisClient())
	protected := httpServer.GetInstance().Group("/groups", authMiddleware)
	protected.Post("/", groupHandler.CreateGroup)
	protected.Get("/my", groupHandler.GetMyGroups)
	protected.Get("/:id", groupHandler.GetGroupByID)
	protected.Get("/", groupHandler.GetAllGroups)
	protected.Put("/:id", groupHandler.UpdateGroup)
	protected.Delete("/:id", groupHandler.DeleteGroup)
}
