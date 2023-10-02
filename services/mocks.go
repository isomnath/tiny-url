package services

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/isomnath/tiny-url/models"
)

/*
Repositories Mocks
*/
type seriesCounterRepoMock struct {
	mock.Mock
}

func (mock *seriesCounterRepoMock) Next(ctx context.Context) (models.Series, error) {
	args := mock.Called(ctx)
	return args.Get(0).(models.Series), args.Error(1)
}

type urlStoreRepoMock struct {
	mock.Mock
}

func (mock *urlStoreRepoMock) InsertMap(ctx context.Context, model models.URLMap) error {
	args := mock.Called(ctx, model)
	return args.Error(0)
}

func (mock *urlStoreRepoMock) Fetch(ctx context.Context, shortenedRep string) (models.URLMap, error) {
	args := mock.Called(ctx, shortenedRep)
	return args.Get(0).(models.URLMap), args.Error(1)
}

type hashStoreRepoMock struct {
	mock.Mock
}

func (mock *hashStoreRepoMock) InsertMap(ctx context.Context, model models.HashMap) error {
	args := mock.Called(ctx, model)
	return args.Error(0)
}

func (mock *hashStoreRepoMock) Fetch(ctx context.Context, hashKey string) (models.HashMap, error) {
	args := mock.Called(ctx, hashKey)
	return args.Get(0).(models.HashMap), args.Error(1)
}

type metricsStoreRepoMock struct {
	mock.Mock
}

func (mock *metricsStoreRepoMock) FetchHighestTransformationSet(ctx context.Context, pageSize int64) ([]models.DomainTransformationCounter, error) {
	args := mock.Called(ctx, pageSize)
	return args.Get(0).([]models.DomainTransformationCounter), args.Error(1)
}

func (mock *metricsStoreRepoMock) FetchHighestTrafficSet(ctx context.Context, pageSize int64) ([]models.DomainRedirectionCounter, error) {
	args := mock.Called(ctx, pageSize)
	return args.Get(0).([]models.DomainRedirectionCounter), args.Error(1)
}

/*
Processors Mock
*/
type sha512ProcessorMock struct {
	mock.Mock
}

func (mock *sha512ProcessorMock) Encode(ctx context.Context, url string) string {
	args := mock.Called(ctx, url)
	return args.String(0)
}

type base62ProcessorMock struct {
	mock.Mock
}

func (mock *base62ProcessorMock) Encode(ctx context.Context, series uint64) string {
	args := mock.Called(ctx, series)
	return args.String(0)
}

type urlProcessorMock struct {
	mock.Mock
}

func (mock *urlProcessorMock) ExtractDomain(ctx context.Context, rawURL string) string {
	args := mock.Called(ctx, rawURL)
	return args.String(0)
}

/*
Metrics Worker Mock
*/
type metricsWorkerMock struct {
	mock.Mock
}

func (mock *metricsWorkerMock) IncrementGenerationCounter(ctx context.Context, domain string) {
	mock.Called(ctx, domain)
	return
}
func (mock *metricsWorkerMock) IncrementRedirectionTrafficCounter(ctx context.Context, domain string) {
	mock.Called(ctx, domain)
	return
}
