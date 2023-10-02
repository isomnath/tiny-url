package services

import (
	"context"

	"github.com/isomnath/tiny-url/errs"
)

type URLRedirectionService struct {
	urlStoreRepo  urlStoreRepo
	urlProcessor  urlProcessor
	metricsWorker metricsWorker
}

func (service *URLRedirectionService) Redirect(ctx context.Context, shortenedRep string) (string, error) {
	urlMap, err := service.urlStoreRepo.Fetch(ctx, shortenedRep)
	if err != nil {
		if err == errs.ErrURLFetchFailureKeyDoesNotExist {
			return "", errs.ErrURLDoesNotExist
		}
		return "", errs.ErrFailedToFetchURLFromURLStore
	}
	domain := service.urlProcessor.ExtractDomain(ctx, urlMap.OriginalURL)
	service.metricsWorker.IncrementRedirectionTrafficCounter(ctx, domain)
	return urlMap.OriginalURL, nil
}

func NewURLRedirectionService(urlStoreRepo urlStoreRepo, urlProcessor urlProcessor,
	metricsWorker metricsWorker) *URLRedirectionService {
	return &URLRedirectionService{
		urlStoreRepo:  urlStoreRepo,
		urlProcessor:  urlProcessor,
		metricsWorker: metricsWorker,
	}
}
