package main

import (
	"log"
	"os"
	"strconv"

	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal"
	"github.com/haserta98/go-rest/internal/models"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/haserta98/go-rest/internal/ws"
	"github.com/haserta98/go-rest/routes"
)

func main() {
	httpCH := make(chan *cmd.HTTPServerImpl)
	dbCH := make(chan *cmd.DB)

	go initHTTPServer(httpCH)
	go initDB(dbCH)

	httpServer := <-httpCH
	db := <-dbCH

	redisClient := internal.NewRedisClient()
	cluster := cmd.NewCluster(redisClient)
	cluster.SendHeartbeat()

	manager := ws.NewWsManager(redisClient, cluster)

	go manager.ListenRedis()

	wsGateway := ws.NewWsGateway(httpServer, manager)
	wsGateway.HandleWebSocket()
	wsGateway.Start()

	ctx := cmd.NewAppContext(httpServer, db)
	if ctx.GetHTTPServer() == nil {
		log.Fatal("HTTP server instance in context is nil")
	}
	if ctx.GetDB() == nil {
		log.Fatal("DB instance in context is nil")
	}

	ctx.GetDB().Migrate(&models.User{})
	ctx.GetDB().Migrate(&models.Group{})
	ctx.GetDB().Migrate(&models.GroupMember{})
	ctx.GetDB().Migrate(&models.Message{})

	ctx.RegisterRepository("User", *repository.NewUserRepository(ctx.GetDB()))
	ctx.RegisterRepository("Group", *repository.NewGroupRepository(ctx.GetDB()))

	routes.InitRoutes(ctx)

	log.Println("Http server başlatıldı")
	log.Println("Application has been started")
	if err := ctx.GetHTTPServer().Listen(); err != nil {
		log.Fatalf("HTTP server başlatılamadı: %v", err)
	}

	select {}
}

func initHTTPServer(ch chan *cmd.HTTPServerImpl) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	parsedPort, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}
	httpServer := cmd.NewHTTPServer(int(parsedPort))
	if (httpServer.GetInstance()) == nil {
		panic("HTTP server instance is nil")
	}
	ch <- httpServer
}

func initDB(ch chan *cmd.DB) {
	db := cmd.NewDB()
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	ch <- db
}
