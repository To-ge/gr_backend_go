package usecase

import (
	"fmt"
	"log"

	"github.com/To-ge/gr_backend_go/domain/entity"
	"github.com/To-ge/gr_backend_go/domain/repository"
	"github.com/To-ge/gr_backend_go/usecase/model"
)

type ITelemetryLogUsecase interface {
	GetTelemetryLogs(input *model.GetTelemetryLogsInput, isAdmin bool) (*model.GetTelemetryLogsOutput, error)
	ToggleTelemetryLogVisibility(input *model.ToggleTelemetryLogVisibilityInput) error
}

type telemetryLogUsecase struct {
	repo repository.ITelemetryLogRepository
}

func NewTelemetryLogUsecase(tlr repository.ITelemetryLogRepository) ITelemetryLogUsecase {
	return &telemetryLogUsecase{
		repo: tlr,
	}
}

func (tlu *telemetryLogUsecase) GetTelemetryLogs(input *model.GetTelemetryLogsInput, isAdmin bool) (*model.GetTelemetryLogsOutput, error) {
	var logs []entity.TelemetryLog
	var err error

	if isAdmin {
		logs, err = tlu.repo.GetTelemetryLogs()
	} else {
		logs, err = tlu.repo.GetPublicTelemetryLogs()
	}
	if err != nil {
		return nil, err
	}

	var list []model.TelemetryLog
	for _, v := range logs {
		list = append(list, model.TelemetryLog{
			ID:            v.ID,
			StartTime:     v.StartTime,
			EndTime:       v.EndTime,
			LocationCount: v.LocationCount,
			IsPublic:      v.IsPublic,
		})
	}

	return &model.GetTelemetryLogsOutput{
		IsPublic: !isAdmin,
		Logs:     list,
	}, nil
}

func (tlu *telemetryLogUsecase) ToggleTelemetryLogVisibility(input *model.ToggleTelemetryLogVisibilityInput) error {
	if err := tlu.repo.ToggleTelemetryLogVisibility(input.Id, input.Visible); err != nil {
		log.Printf("Error: telemetryLogUsecase ID(%d), %v\n", input.Id, err)
		return fmt.Errorf("an error occurred in the database")
	}
	return nil
}
