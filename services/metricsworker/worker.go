package metricsworker

import (
	"context"

	"github.com/isomnath/tiny-url/log"
)

type MetricsWorker struct {
	metricsStoreRepo metricsStoreRepo
}

func (worker *MetricsWorker) IncrementGenerationCounter(ctx context.Context, domain string) {
	go func() {
		err := worker.metricsStoreRepo.IncrementTransformationCounter(ctx, domain)
		if err != nil {
			log.Log.WithError(err).Warnf(err.Error())
		}
	}()
}

func (worker *MetricsWorker) IncrementRedirectionTrafficCounter(ctx context.Context, domain string) {
	go func() {
		err := worker.metricsStoreRepo.IncrementRedirectionTrafficCounter(ctx, domain)
		if err != nil {
			log.Log.WithError(err).Warnf(err.Error())
		}
	}()
}

func NewMetricsWorker(metricsStoreRepo metricsStoreRepo) *MetricsWorker {
	return &MetricsWorker{
		metricsStoreRepo: metricsStoreRepo,
	}
}
