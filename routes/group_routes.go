package routes

import (
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal/handler"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/haserta98/go-rest/internal/service"
)

func initGroupEndpoints(appCtx *cmd.AppContext) {
	httpServer := appCtx.GetHTTPServer()
	groupService := service.NewGroupService(appCtx.GetRepository("Group").(repository.GroupRepository))
	groupHandler := handler.NewGroupHandler(groupService)

	httpServer.GetInstance().Post("/groups", groupHandler.CreateGroup)
	httpServer.GetInstance().Get("/groups/:id", groupHandler.GetGroupByID)
	httpServer.GetInstance().Get("/groups", groupHandler.GetAllGroups)
	httpServer.GetInstance().Put("/groups/:id", groupHandler.UpdateGroup)
	httpServer.GetInstance().Delete("/groups/:id", groupHandler.DeleteGroup)
}
