package handler

import (
	"io"
	"log"
	"math"
	"time"

	v1 "github.com/To-ge/gr_backend_go/adapter/grpc/api/gen/go/v1"
	"github.com/To-ge/gr_backend_go/pkg"
	"github.com/To-ge/gr_backend_go/usecase"
	"github.com/To-ge/gr_backend_go/usecase/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type telemetryHandler struct {
	v1.UnimplementedTelemetryServiceServer
	usecase usecase.ITelemetryUsecase
}

func NewTelemetryHandler(tu usecase.ITelemetryUsecase) *telemetryHandler {
	return &telemetryHandler{
		usecase: tu,
	}
}

func (th *telemetryHandler) SendLocation(stream v1.TelemetryService_SendLocationServer) error {
	log.Println("telemetryHandler.SendLocation started.")
	defer log.Println("telemetryHandler.SendLocation ended.")
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			log.Println("telemetryHandler.SendLocation client is done.")
			th.usecase.Stop()
			return nil
		default:
			req, err := stream.Recv()
			currentTime := float64(time.Now().UnixMicro()) / math.Pow10(6)
			pkg.InputLocationLogger.Printf(",%f,%v,%v,%v\n", currentTime, req.GetLatitude(), req.GetLongitude(), req.GetAltitude())
			if err == io.EOF {
				log.Println("telemetryHandler.SendLocation Client closed the stream")
				th.usecase.Stop()
				return nil
			} else if err != nil {
				log.Printf("telemetryHandler.SendLocation error: %s", err.Error())
				th.usecase.Stop()
				return status.Errorf(codes.Internal, err.Error())
			}

			input := &model.SendLocationInput{
				Location: model.Location{
					Timestamp: req.Timestamp,
					Latitude:  req.Latitude,
					Longitude: req.Longitude,
					Altitude:  float64(req.Altitude),
				},
			}
			if _, err := th.usecase.SendLocation(input); err != nil {
				log.Printf("telemetryHandler.usecase.SendLocation error: %s", err.Error())
				continue
			}
		}
	}
}
