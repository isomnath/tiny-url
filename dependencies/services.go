package dependencies

import (
	"github.com/isomnath/tiny-url/services"
	"github.com/isomnath/tiny-url/services/metricsworker"
	"github.com/isomnath/tiny-url/services/processors"
)

type Services struct {
	AnalyticsService      *services.AnalyticsService
	URLGenerationService  *services.URLGenerationService
	URLRedirectionService *services.URLRedirectionService
}

func initializeServices(repositories *Repositories) *Services {
	base62Processor := processors.NewBase62Processor()
	sha512Processor := processors.NewSha512Processor()
	urlProcessor := processors.NewURLProcessor()
	metricsWorker := metricsworker.NewMetricsWorker(repositories.MetricsStoreRepo)
	return &Services{
		AnalyticsService:      services.NewAnalyticsService(repositories.MetricsStoreRepo),
		URLRedirectionService: services.NewURLRedirectionService(repositories.URLStoreRepo, urlProcessor, metricsWorker),
		URLGenerationService: services.NewURLGenerationService(repositories.SeriesCounterRepo, repositories.URLStoreRepo,
			repositories.HashStoreRepo, sha512Processor, base62Processor, urlProcessor, metricsWorker),
	}
}
