package main

import (
	"fmt"
	"log"
	"log-service/contracts"
	"log-service/data/adaptor"
	"log-service/data/repository"
	"net"
	"net/http"
	"net/rpc"
	"os"

	logsvc "log-service/logSvc"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	rpcPort  = "5001"
	grpcPort = "50001"
	webPort  = "80"
)

type Config struct {
	mongo adaptor.MongoConfig
}

type HTTPServer struct {
	svc *logsvc.Service
}

func main() {
	conf := loadConfig()
	client, err := adaptor.ConnectToMongo(conf.mongo)
	if err != nil {
		log.Panic(err)
	}
	// Close MongoDB connection
	defer adaptor.Disconnect(client)

	logRepo := repository.New(client)
	logSvc := logsvc.New(logRepo)

	app := HTTPServer{
		svc: logSvc,
	}

	// Register RPC Server
	rpcServer := &RPCServer{
		svc: logSvc,
	}
	rpc.Register(rpcServer)
	go rpcServer.serveRPC()

	// Serve gRPC Server
	grpcServer := &GRPCLogServer{
		svc: logSvc,
	}
	go serveGRPC(grpcServer)

	// start webserver
	app.serve()
}

func (app *HTTPServer) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.SetRoutes(),
	}

	log.Printf("The server is now running on %s Address.", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Because of the following error, server had to stopped: %s", err)
	}
}

func (r *RPCServer) serveRPC() {
	RPCAddress := fmt.Sprintf("0.0.0.0:%s", rpcPort)
	listener, lErr := net.Listen("tcp", RPCAddress)
	if lErr != nil {
		log.Fatal("Error listening to RPC address")
	}
	defer listener.Close()

	for {
		conn, connErr := listener.Accept()
		if connErr != nil {
			log.Fatal("Error connecting/accepting to RPC Server")
		}
		go rpc.ServeConn(conn)

	}
}

func serveGRPC(grpcLogServer *GRPCLogServer) {
	gRPCAddress := fmt.Sprintf("0.0.0.0:%s", grpcPort)
	listener, lErr := net.Listen("tcp", gRPCAddress)
	if lErr != nil {
		log.Fatal("Error listening to gRPC address")
	}
	defer listener.Close()

	server := grpc.NewServer()

	contracts.RegisterLogServiceServer(server, grpcLogServer)
	log.Println("gRPC server start listening on %s port", grpcPort)

	if sErr := server.Serve(listener); sErr != nil {
		log.Fatal("Error serving gRPC server: ", sErr)
	}
}

func loadConfig() Config {
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		log.Fatal("Error loading .env file", err)
	}

	return Config{
		mongo: adaptor.MongoConfig{
			Username:     os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
			Password:     os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
			MongoAddress: os.Getenv("MONGO_ADDRESS"),
			DB:           os.Getenv("MONGO_INITDB_DATABASE"),
		},
	}
}
