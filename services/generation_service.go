package services

import (
	"context"
	"fmt"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/contracts"
	"github.com/isomnath/tiny-url/errs"
	"github.com/isomnath/tiny-url/models"
)

type URLGenerationService struct {
	seriesCounterRepo seriesCounterRepo
	urlStoreRepo      urlStoreRepo
	hashStoreRepo     hashStoreRepo
	sha512Processor   sha512Processor
	base62Processor   base62Processor
	urlProcessor      urlProcessor
	metricsWorker     metricsWorker
}

func (service *URLGenerationService) Generate(ctx context.Context, request contracts.GenerateRequest) (contracts.GenerateResponse, error) {
	urlHash := service.sha512Processor.Encode(ctx, request.OriginalURL)
	hashMap, err := service.hashStoreRepo.Fetch(ctx, urlHash)
	if err != nil {
		if err == errs.ErrHashFetchFailureKeyDoesNotExist {
			series, err := service.seriesCounterRepo.Next(ctx)
			if err != nil {
				return contracts.GenerateResponse{}, errs.ErrFailedToFetchNextValueInSeries
			}
			shortenedRep := service.base62Processor.Encode(ctx, uint64(series))
			err = service.hashStoreRepo.InsertMap(ctx, models.HashMap{HashKey: urlHash, ShortenedRep: shortenedRep})
			if err != nil {
				return contracts.GenerateResponse{}, errs.ErrFailedToSaveHashInStore
			}
			err = service.urlStoreRepo.InsertMap(ctx, models.URLMap{ShortenedRep: shortenedRep, OriginalURL: request.OriginalURL})
			if err != nil {
				return contracts.GenerateResponse{}, errs.ErrFailedToSaveURLInStore
			}
			domain := service.urlProcessor.ExtractDomain(ctx, request.OriginalURL)
			service.metricsWorker.IncrementGenerationCounter(ctx, domain)
			return contracts.GenerateResponse{TinyURL: fmt.Sprintf("%s/%s", config.GetAppDNS(), shortenedRep)}, nil
		}
		return contracts.GenerateResponse{}, errs.ErrFailedToFetchHashFromHashStore
	}
	urlMap, err := service.urlStoreRepo.Fetch(ctx, hashMap.ShortenedRep)
	if err != nil {
		return contracts.GenerateResponse{}, errs.ErrFailedToFetchURLFromURLStore
	}

	if urlMap.OriginalURL != request.OriginalURL {
		return contracts.GenerateResponse{}, errs.ErrURLGenerationHashCollisionURL
	}

	return contracts.GenerateResponse{TinyURL: fmt.Sprintf("%s/%s", config.GetAppDNS(), urlMap.ShortenedRep)}, nil
}

func NewURLGenerationService(seriesCounterRepo seriesCounterRepo, urlStoreRepo urlStoreRepo,
	hashStoreRepo hashStoreRepo, sha512Processor sha512Processor, base62Processor base62Processor, urlProcessor urlProcessor,
	metricsWorker metricsWorker) *URLGenerationService {
	return &URLGenerationService{
		seriesCounterRepo: seriesCounterRepo,
		urlStoreRepo:      urlStoreRepo,
		hashStoreRepo:     hashStoreRepo,
		sha512Processor:   sha512Processor,
		base62Processor:   base62Processor,
		urlProcessor:      urlProcessor,
		metricsWorker:     metricsWorker,
	}
}
