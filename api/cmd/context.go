package cmd

import (
	"log"

	"github.com/haserta98/go-rest/internal"
)

type AppContext struct {
	httpServer   *HTTPServerImpl
	db           *DB
	redisClient  *internal.RedisClient
	repositories map[string]interface{}
	services     map[string]interface{}
}

func NewAppContext(httpServer *HTTPServerImpl, db *DB, redisClient *internal.RedisClient) *AppContext {
	return &AppContext{
		httpServer:   httpServer,
		db:           db,
		redisClient:  redisClient,
		repositories: make(map[string]interface{}),
		services:     make(map[string]interface{}),
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

func (ctx *AppContext) RegisterService(name string, service interface{}) {
	if name == "" {
		panic("Service name cannot be empty")
	}
	if service == nil {
		panic("Service instance cannot be nil")
	}
	if _, exists := ctx.services[name]; exists {
		panic("Service with this name already exists: " + name)
	}
	log.Printf("Registering service: %s", name)
	ctx.services[name] = service
}

func (ctx *AppContext) GetService(name string) interface{} {
	return ctx.services[name]
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

func (ctx *AppContext) GetRedisClient() *internal.RedisClient {
	return ctx.redisClient
}
