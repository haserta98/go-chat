package routes

import (
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal/handler"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/haserta98/go-rest/internal/service"
)

func initUserEndpoints(appCtx *cmd.AppContext) {

	httpServer := appCtx.GetHTTPServer()
	userService := service.NewUserService(appCtx.GetRepository("User").(repository.UserRepository))
	userHandler := handler.NewUserHandler(userService)

	httpServer.GetInstance().Post("/users", userHandler.CreateUser)
	httpServer.GetInstance().Get("/users/:id", userHandler.GetUserByID)
	httpServer.GetInstance().Get("/users", userHandler.GetAllUsers)
	httpServer.GetInstance().Put("/users/:id", userHandler.UpdateUser)
	httpServer.GetInstance().Delete("/users/:id", userHandler.DeleteUser)
}
