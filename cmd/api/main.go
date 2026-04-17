package main

import (
	"log"
	"os"
	"strconv"

	"github.com/haserta98/go-rest/cmd"
	"github.com/haserta98/go-rest/internal"
	"github.com/haserta98/go-rest/internal/models"
	"github.com/haserta98/go-rest/internal/repository"
	"github.com/haserta98/go-rest/internal/service"
	"github.com/haserta98/go-rest/internal/ws"
	"github.com/haserta98/go-rest/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Fatal("NODE_ID environment variable is not set")
	}

	httpCH := make(chan *cmd.HTTPServerImpl)
	dbCH := make(chan *cmd.DB)

	go initHTTPServer(httpCH)
	go initDB(dbCH)

	httpServer := <-httpCH
	db := <-dbCH

	redisClient := internal.NewRedisClient()
	cluster := cmd.NewCluster(redisClient, nodeID)
	cluster.SendHeartbeat()

	manager := ws.NewWsManager(redisClient, cluster)
	manager.Start()

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
	ctx.RegisterRepository("Message", repository.NewMessageRepository(ctx.GetDB()))

	messageService := service.NewMessageService(wsGateway, ctx.GetRepository("Message").(*repository.MessageRepository))
	messageService.RegisterEventHandlers()

	ctx.RegisterService("MessageService", messageService)
	ctx.RegisterService("WsGateway", wsGateway)

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
