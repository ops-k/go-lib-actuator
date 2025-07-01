package services

import (
	"context"

	models "github.com/ops-k/go-lib-actuator/models"
)

// PingService struct
type PingService struct {
}

func NewPingService() *PingService {
	return &PingService{}
}
func (svc *PingService) Ping(ctx context.Context) *models.PingResponse {
	return &models.PingResponse{
		Message: "pong",
	}
}
