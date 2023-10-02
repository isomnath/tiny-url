package services

import (
	"context"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/contracts"
)

type AnalyticsService struct {
	metricsStoreRepo metricsStoreRepo
}

func (service *AnalyticsService) FetchTopDomainsByTransformations(ctx context.Context) (contracts.TransformationResponse, error) {
	transformations, err := service.metricsStoreRepo.FetchHighestTransformationSet(ctx, int64(config.GetAnalyticsConfigs().TopTransformationsPageSize()))
	if err != nil {
		return contracts.TransformationResponse{}, err
	}

	var topTransformations []contracts.Transformation
	for _, transformation := range transformations {
		topTransformations = append(topTransformations, contracts.Transformation{
			Domain: transformation.Domain,
			Count:  transformation.Transformations,
		})
	}

	return contracts.TransformationResponse{TopTransformations: topTransformations}, nil
}

func (service *AnalyticsService) FetchTopDomainsByTraffic(ctx context.Context) (contracts.TrafficResponse, error) {
	traffic, err := service.metricsStoreRepo.FetchHighestTrafficSet(ctx, int64(config.GetAnalyticsConfigs().TopTrafficPageSize()))
	if err != nil {
		return contracts.TrafficResponse{}, err
	}

	var topTraffic []contracts.Traffic
	for _, trf := range traffic {
		topTraffic = append(topTraffic, contracts.Traffic{
			Domain: trf.Domain,
			Count:  trf.Traffic,
		})
	}

	return contracts.TrafficResponse{TopTraffic: topTraffic}, nil
}

func NewAnalyticsService(metricsStoreRepo metricsStoreRepo) *AnalyticsService {
	return &AnalyticsService{
		metricsStoreRepo: metricsStoreRepo,
	}
}
