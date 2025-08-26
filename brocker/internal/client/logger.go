package client

import (
	"brocker/internal/dto"
	"log"
	"net/rpc"
)

type LoggerRPCClient struct {
	client *rpc.Client
}

func NewLoggerRPCClient(addr string) (*LoggerRPCClient, error) {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &LoggerRPCClient{client: client}, nil
}

// Log sends a log entry to the logger service via RPC.
func (c *LoggerRPCClient) Log(payload dto.LogRPCRequestPayload) (string, error) {
	var result string

	err := c.client.Call("RPCServer.LogInfo", payload, &result)
	if err != nil {
		log.Printf("ERROR: RPC call to logger failed: %v", err)
		return "", err
	}

	log.Printf("Logger service replied: %s", result)
	return result, nil
}