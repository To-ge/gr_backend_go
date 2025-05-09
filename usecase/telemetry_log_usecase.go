package usecase

import (
	"github.com/To-ge/gr_backend_go/domain/repository"
	"github.com/To-ge/gr_backend_go/usecase/model"
)

type ITelemetryLogUsecase interface {
	GetTelemetryLogs(*model.GetTelemetryLogsInput) (*model.GetTelemetryLogsOutput, error)
}

type telemetryLogUsecase struct {
	repo repository.ITelemetryLogRepository
}

func NewTelemetryLogUsecase(tlr repository.ITelemetryLogRepository) ITelemetryLogUsecase {
	return &telemetryLogUsecase{
		repo: tlr,
	}
}

func (tlu *telemetryLogUsecase) GetTelemetryLogs(input *model.GetTelemetryLogsInput) (*model.GetTelemetryLogsOutput, error) {
	logs, err := tlu.repo.GetPublicTelemetryLogs()
	if err != nil {
		return nil, err
	}

	var list []model.TelemetryLog
	for _, v := range logs {
		list = append(list, model.TelemetryLog{
			StartTime:     v.StartTime,
			EndTime:       v.EndTime,
			LocationCount: v.LocationCount,
		})
	}

	return &model.GetTelemetryLogsOutput{
		Logs: list,
	}, nil
}
