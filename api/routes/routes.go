package routes

import (
	"github.com/haserta98/go-rest/cmd"
)

func InitRoutes(ctx *cmd.AppContext) {
	initUserEndpoints(ctx)
	initGroupEndpoints(ctx)
	initMessageEndpoints(ctx)
}
