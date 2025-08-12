package logsvc

import "log-service/domain"

type LogRepo interface {
	Insert(dlog domain.LogEntry) error
	All() ([]*domain.LogEntry, error)
}