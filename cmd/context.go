package cmd

import (
	"log"
)

type AppContext struct {
	httpServer   *HTTPServerImpl
	db           *DB
	repositories map[string]interface{}
}

func NewAppContext(httpServer *HTTPServerImpl, db *DB) *AppContext {
	return &AppContext{
		httpServer:   httpServer,
		db:           db,
		repositories: make(map[string]interface{}),
	}
}

func (ctx *AppContext) RegisterRepository(name string, repo interface{}) {
	if name == "" {
		panic("Repository name cannot be empty")
	}
	if repo == nil {
		panic("Repository instance cannot be nil")
	}
	if _, exists := ctx.repositories[name]; exists {
		panic("Repository with this name already exists: " + name)
	}
	log.Printf("Registering repository: %s", name)
	ctx.repositories[name] = repo
}

func (ctx *AppContext) GetRepository(name string) interface{} {
	return ctx.repositories[name]
}

func (ctx *AppContext) GetHTTPServer() *HTTPServerImpl {
	return ctx.httpServer
}

func (ctx *AppContext) GetDB() *DB {
	return ctx.db
}
