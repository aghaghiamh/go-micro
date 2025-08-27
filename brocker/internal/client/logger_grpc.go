package client

import (
	"brocker/internal/contracts"
	"brocker/internal/dto"
	"context"

	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoggerGRPCClient struct {
	client contracts.LogServiceClient
}

func NewLoggerGRPCClient(grpcURI string) (*LoggerGRPCClient, error) {
	conn, cErr := grpc.NewClient(grpcURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if cErr != nil {
		errMsg := fmt.Sprintf("Error: couldn't conenct to the grpc server: ", cErr)
		log.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	logClient := contracts.NewLogServiceClient(conn)
	return &LoggerGRPCClient{
		client: logClient,
	}, nil
}

func (c *LoggerGRPCClient) WriteLog(ctx context.Context, logPayload dto.LogPayload) (string, error) {
	req := &contracts.LogRequest{
		LogEntry: &contracts.Log{
			Name:  logPayload.Name,
			Data:  logPayload.Data,
			Level: contracts.LogLevel_LOG_LEVEL_INFO, // Example log level
		},
	}

	res, err := c.client.WriteLog(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Result, nil
}
