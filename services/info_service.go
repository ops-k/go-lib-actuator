package services

import (
	"context"

	models "github.com/ops-k/go-lib-actuator/models"
)

type InfoContributor interface {
	GetKey() string
	GetValue() interface{}
}

type InfoContributors []InfoContributor

// ActuatorInfoService struct
type ActuatorInfoService struct {
	infoContributors InfoContributors
}

func NewActuatorInfoService(infoContributors InfoContributors) *ActuatorInfoService {
	return &ActuatorInfoService{
		infoContributors: infoContributors,
	}
}

func (svc *ActuatorInfoService) GetInfo(ctx context.Context) map[string]interface{} {
	result := make(map[string]interface{})
	for _, infoContributor := range svc.infoContributors {
		result[infoContributor.GetKey()] = infoContributor.GetValue()
	}
	return result
}

type BuildInfoContributor struct {
	buildInfo *models.BuildInfo
}

func NewBuildInfoContributor(buildInfo *models.BuildInfo) *BuildInfoContributor {
	return &BuildInfoContributor{
		buildInfo: buildInfo,
	}
}

func (o *BuildInfoContributor) GetKey() string {
	return "build"
}

func (o *BuildInfoContributor) GetValue() interface{} {
	return o.buildInfo
}
