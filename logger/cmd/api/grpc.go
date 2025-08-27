package main

import (
	"context"
	"fmt"
	"log-service/contracts"
	"log-service/domain"
	logsvc "log-service/logSvc"
)

type GRPCLogServer struct {
	contracts.UnimplementedLogServiceServer
	svc *logsvc.Service
}

func (r *GRPCLogServer) WriteLog(ctx context.Context, req *contracts.LogRequest) (*contracts.LogResponse, error) {
	input := req.GetLogEntry()
	logEntry := domain.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}
	if svcErr := r.svc.WriteLog(logEntry); svcErr != nil {
		res := &contracts.LogResponse{Result: "failed"}
		return res, fmt.Errorf("handler error: calling svc.WriteLog: %w", svcErr)
	}

	res := &contracts.LogResponse{
		Result: "Logged using gRPC server",
	}
	return res, nil
}
