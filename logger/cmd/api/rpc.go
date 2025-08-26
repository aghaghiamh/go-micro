package main

import (
	"fmt"
	"log-service/domain"
	logsvc "log-service/logSvc"
)

type RPCServer struct {
	svc *logsvc.Service
}

func (r *RPCServer) LogInfo(payload WriteLogRPCRequest, resp *string) error {
	logEntry := domain.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	}
	if svcErr := r.svc.WriteLog(logEntry); svcErr != nil {
		return fmt.Errorf("handler error: calling svc.WriteLog: %w", svcErr)
	}
	
	*resp = fmt.Sprintf("Proccessed LogInfo rpc call with %s name", logEntry.Name)
	return nil
}