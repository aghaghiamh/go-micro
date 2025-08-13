package logsvc

import "log-service/domain"

type Service struct {
	repo LogRepo
}

func New(repo LogRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (svc Service) WriteLog(log domain.LogEntry) error {
	iErr := svc.repo.Insert(log)
	return iErr
}
