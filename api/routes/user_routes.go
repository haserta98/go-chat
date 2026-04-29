package routes

import (
	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal/handler"
	"github.com/haserta98/go-rest/internal/middleware"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/haserta98/go-rest/internal/service"
)

func initUserEndpoints(appCtx *cmd.AppContext) {

	httpServer := appCtx.GetHTTPServer()
	userService := service.NewUserService(appCtx.GetRepository("User").(*repository.UserRepository))
	userHandler := handler.NewUserHandler(userService, appCtx.GetRedisClient())

	httpServer.GetInstance().Post("/register", userHandler.CreateUser)
	httpServer.GetInstance().Post("/login", userHandler.LoginUser)
	httpServer.GetInstance().Post("/logout", userHandler.LogoutUser)

	// Protected routes
	authMiddleware := middleware.NewAuthMiddleware(appCtx.GetRedisClient())
	protected := httpServer.GetInstance().Group("/users", authMiddleware)
	protected.Get("/me", userHandler.GetMe)
	protected.Get("/contacts", userHandler.GetContacts)
	protected.Get("/contacts/online", userHandler.GetContactsOnlineStatus)
	protected.Post("/contacts", userHandler.AddContact)
	protected.Delete("/contacts", userHandler.RemoveContact)
	protected.Get("/", userHandler.GetAllUsers)

	// Wildcard routes MUST come after static paths
	protected.Get("/:id", userHandler.GetUserByID)
	protected.Put("/:id", userHandler.UpdateUser)
	protected.Delete("/:id", userHandler.DeleteUser)
}
