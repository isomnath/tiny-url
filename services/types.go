package services

import (
	"context"

	"github.com/isomnath/tiny-url/models"
)

/*
Repositories
*/
type seriesCounterRepo interface {
	Next(ctx context.Context) (models.Series, error)
}

type urlStoreRepo interface {
	InsertMap(ctx context.Context, model models.URLMap) error
	Fetch(ctx context.Context, shortenedRep string) (models.URLMap, error)
}

type hashStoreRepo interface {
	InsertMap(ctx context.Context, model models.HashMap) error
	Fetch(ctx context.Context, hashKey string) (models.HashMap, error)
}

type metricsStoreRepo interface {
	FetchHighestTransformationSet(ctx context.Context, pageSize int64) ([]models.DomainTransformationCounter, error)
	FetchHighestTrafficSet(ctx context.Context, pageSize int64) ([]models.DomainRedirectionCounter, error)
}

/*
Processors
*/
type sha512Processor interface {
	Encode(ctx context.Context, url string) string
}

type base62Processor interface {
	Encode(ctx context.Context, series uint64) string
}

type urlProcessor interface {
	ExtractDomain(ctx context.Context, rawURL string) string
}

/*
Metrics Worker
*/
type metricsWorker interface {
	IncrementGenerationCounter(ctx context.Context, domain string)
	IncrementRedirectionTrafficCounter(ctx context.Context, domain string)
}
