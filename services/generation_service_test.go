package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/contracts"
	"github.com/isomnath/tiny-url/errs"
	"github.com/isomnath/tiny-url/log"
	"github.com/isomnath/tiny-url/models"
)

type URLGenerationServiceTestSuite struct {
	suite.Suite
	ctx               context.Context
	seriesCounterRepo *seriesCounterRepoMock
	urlStoreRepo      *urlStoreRepoMock
	hashStoreRepo     *hashStoreRepoMock
	sha512Processor   *sha512ProcessorMock
	base62Processor   *base62ProcessorMock
	urlProcessor      *urlProcessorMock
	metricsWorker     *metricsWorkerMock
	service           *URLGenerationService
}

func (suite *URLGenerationServiceTestSuite) SetupTest() {
	config.LoadBaseConfig()
	log.Setup()
	suite.ctx = context.Background()
	suite.seriesCounterRepo = new(seriesCounterRepoMock)
	suite.urlStoreRepo = new(urlStoreRepoMock)
	suite.hashStoreRepo = new(hashStoreRepoMock)
	suite.sha512Processor = new(sha512ProcessorMock)
	suite.base62Processor = new(base62ProcessorMock)
	suite.urlProcessor = new(urlProcessorMock)
	suite.metricsWorker = new(metricsWorkerMock)
	suite.service = NewURLGenerationService(suite.seriesCounterRepo, suite.urlStoreRepo, suite.hashStoreRepo,
		suite.sha512Processor, suite.base62Processor, suite.urlProcessor, suite.metricsWorker)
}

func (suite *URLGenerationServiceTestSuite) TestGenerateWhenHashDoesNotExistShouldReturnErrorWhenFailedToFetchNextSeriesValue() {
	url := "https://www.google.com?q=infracloud"
	urlHash := "559ae209f23cf3866edd0b6bf6c11c8baa128cfc620736d105d582ff8082d1baca509d100daafc165e1ff5d81f1bd3535c9cfa26dba96569a9bec318aee6e756"
	suite.sha512Processor.On("Encode", suite.ctx, url).Return(urlHash)
	suite.hashStoreRepo.On("Fetch", suite.ctx, urlHash).Return(models.HashMap{}, errs.ErrHashFetchFailureKeyDoesNotExist)
	suite.urlProcessor.On("ExtractDomain", suite.ctx, url).Return("google")
	suite.seriesCounterRepo.On("Next", suite.ctx).Return(models.Series(0), errs.ErrSeriesCounterFailure)

	response, err := suite.service.Generate(suite.ctx, contracts.GenerateRequest{OriginalURL: url})
	suite.Equal(errs.ErrFailedToFetchNextValueInSeries, err)
	suite.Equal(contracts.GenerateResponse{}, response)
}

func (suite *URLGenerationServiceTestSuite) TestGenerateWhenHashDoesNotExistShouldReturnErrorWhenFailedToSaveHashMappingInStore() {
	url := "https://www.google.com?q=infracloud"
	urlHash := "559ae209f23cf3866edd0b6bf6c11c8baa128cfc620736d105d582ff8082d1baca509d100daafc165e1ff5d81f1bd3535c9cfa26dba96569a9bec318aee6e756"
	shortenedRep := "FMJRwWE"
	series := models.Series(1000)
	hashMap := models.HashMap{
		HashKey:      urlHash,
		ShortenedRep: shortenedRep,
	}
	suite.sha512Processor.On("Encode", suite.ctx, url).Return(urlHash)
	suite.hashStoreRepo.On("Fetch", suite.ctx, urlHash).Return(models.HashMap{}, errs.ErrHashFetchFailureKeyDoesNotExist)
	suite.urlProcessor.On("ExtractDomain", suite.ctx, url).Return("google")
	suite.seriesCounterRepo.On("Next", suite.ctx).Return(series, nil)
	suite.base62Processor.On("Encode", suite.ctx, uint64(series)).Return(shortenedRep)
	suite.hashStoreRepo.On("InsertMap", suite.ctx, hashMap).Return(errs.ErrHashInsertFailure)

	response, err := suite.service.Generate(suite.ctx, contracts.GenerateRequest{OriginalURL: url})
	suite.Equal(errs.ErrFailedToSaveHashInStore, err)
	suite.Equal(contracts.GenerateResponse{}, response)
}

func (suite *URLGenerationServiceTestSuite) TestGenerateWhenHashDoesNotExistShouldReturnErrorWhenFailedToSaveURLMappingInStore() {
	url := "https://www.google.com?q=infracloud"
	urlHash := "559ae209f23cf3866edd0b6bf6c11c8baa128cfc620736d105d582ff8082d1baca509d100daafc165e1ff5d81f1bd3535c9cfa26dba96569a9bec318aee6e756"
	shortenedRep := "FMJRwWE"
	series := models.Series(1000)
	hashMap := models.HashMap{
		HashKey:      urlHash,
		ShortenedRep: shortenedRep,
	}
	urlMap := models.URLMap{
		ShortenedRep: shortenedRep,
		OriginalURL:  url,
	}
	suite.sha512Processor.On("Encode", suite.ctx, url).Return(urlHash)
	suite.hashStoreRepo.On("Fetch", suite.ctx, urlHash).Return(models.HashMap{}, errs.ErrHashFetchFailureKeyDoesNotExist)
	suite.urlProcessor.On("ExtractDomain", suite.ctx, url).Return("google")
	suite.seriesCounterRepo.On("Next", suite.ctx).Return(series, nil)
	suite.base62Processor.On("Encode", suite.ctx, uint64(series)).Return(shortenedRep)
	suite.hashStoreRepo.On("InsertMap", suite.ctx, hashMap).Return(nil)
	suite.urlStoreRepo.On("InsertMap", suite.ctx, urlMap).Return(errs.ErrURLInsertFailure)

	response, err := suite.service.Generate(suite.ctx, contracts.GenerateRequest{OriginalURL: url})
	suite.Equal(errs.ErrFailedToSaveURLInStore, err)
	suite.Equal(contracts.GenerateResponse{}, response)
}

func (suite *URLGenerationServiceTestSuite) TestGenerateWhenHashDoesNotExistShouldReturnShortenedRepSuccessfully() {
	url := "https://www.google.com?q=infracloud"
	urlHash := "559ae209f23cf3866edd0b6bf6c11c8baa128cfc620736d105d582ff8082d1baca509d100daafc165e1ff5d81f1bd3535c9cfa26dba96569a9bec318aee6e756"
	shortenedRep := "FMJRwWE"
	series := models.Series(1000)
	hashMap := models.HashMap{
		HashKey:      urlHash,
		ShortenedRep: shortenedRep,
	}
	urlMap := models.URLMap{
		ShortenedRep: shortenedRep,
		OriginalURL:  url,
	}
	suite.sha512Processor.On("Encode", suite.ctx, url).Return(urlHash)
	suite.hashStoreRepo.On("Fetch", suite.ctx, urlHash).Return(models.HashMap{}, errs.ErrHashFetchFailureKeyDoesNotExist)
	suite.urlProcessor.On("ExtractDomain", suite.ctx, url).Return("google")
	suite.seriesCounterRepo.On("Next", suite.ctx).Return(series, nil)
	suite.base62Processor.On("Encode", suite.ctx, uint64(series)).Return(shortenedRep)
	suite.hashStoreRepo.On("InsertMap", suite.ctx, hashMap).Return(nil)
	suite.urlStoreRepo.On("InsertMap", suite.ctx, urlMap).Return(nil)
	suite.metricsWorker.On("IncrementGenerationCounter", suite.ctx, "google")

	response, err := suite.service.Generate(suite.ctx, contracts.GenerateRequest{OriginalURL: url})
	suite.NoError(err)
	suite.Equal(contracts.GenerateResponse{TinyURL: "http://localhost:8181/FMJRwWE"}, response)
}

func (suite *URLGenerationServiceTestSuite) TestGenerateReturnErrorWhenFailedToFetchHashFromHashStore() {
	url := "https://www.google.com?q=infracloud"
	urlHash := "559ae209f23cf3866edd0b6bf6c11c8baa128cfc620736d105d582ff8082d1baca509d100daafc165e1ff5d81f1bd3535c9cfa26dba96569a9bec318aee6e756"
	suite.sha512Processor.On("Encode", suite.ctx, url).Return(urlHash)
	suite.hashStoreRepo.On("Fetch", suite.ctx, urlHash).Return(models.HashMap{}, errs.ErrHashFetchFailure)

	response, err := suite.service.Generate(suite.ctx, contracts.GenerateRequest{OriginalURL: url})
	suite.Equal(errs.ErrFailedToFetchHashFromHashStore, err)
	suite.Equal(contracts.GenerateResponse{}, response)
}

func (suite *URLGenerationServiceTestSuite) TestGenerateWhenHashExistsShouldReturnHashCollisionError() {
	url := "https://www.google.com?q=infracloud"
	urlHash := "559ae209f23cf3866edd0b6bf6c11c8baa128cfc620736d105d582ff8082d1baca509d100daafc165e1ff5d81f1bd3535c9cfa26dba96569a9bec318aee6e756"
	shortenedRep := "FMJRwWE"
	hashMap := models.HashMap{
		HashKey:      urlHash,
		ShortenedRep: shortenedRep,
	}
	urlMap := models.URLMap{
		ShortenedRep: shortenedRep,
		OriginalURL:  "https://www.stackoverflow.com?q=infracloud",
	}
	suite.sha512Processor.On("Encode", suite.ctx, url).Return(urlHash)
	suite.hashStoreRepo.On("Fetch", suite.ctx, urlHash).Return(hashMap, nil)
	suite.urlStoreRepo.On("Fetch", suite.ctx, shortenedRep).Return(urlMap, nil)

	response, err := suite.service.Generate(suite.ctx, contracts.GenerateRequest{OriginalURL: url})
	suite.Equal(errs.ErrURLGenerationHashCollisionURL, err)
	suite.Equal(contracts.GenerateResponse{}, response)
}

func (suite *URLGenerationServiceTestSuite) TestGenerateWhenHashExistsShouldReturnErrorWhenFailedToFetchExistingShortenedRep() {
	url := "https://www.google.com?q=infracloud"
	urlHash := "559ae209f23cf3866edd0b6bf6c11c8baa128cfc620736d105d582ff8082d1baca509d100daafc165e1ff5d81f1bd3535c9cfa26dba96569a9bec318aee6e756"
	shortenedRep := "FMJRwWE"
	hashMap := models.HashMap{
		HashKey:      urlHash,
		ShortenedRep: shortenedRep,
	}
	urlMap := models.URLMap{
		ShortenedRep: shortenedRep,
		OriginalURL:  url,
	}
	suite.sha512Processor.On("Encode", suite.ctx, url).Return(urlHash)
	suite.hashStoreRepo.On("Fetch", suite.ctx, urlHash).Return(hashMap, nil)
	suite.urlStoreRepo.On("Fetch", suite.ctx, shortenedRep).Return(urlMap, errs.ErrURLFetchFailure)

	response, err := suite.service.Generate(suite.ctx, contracts.GenerateRequest{OriginalURL: url})
	suite.Equal(errs.ErrFailedToFetchURLFromURLStore, err)
	suite.Equal(contracts.GenerateResponse{}, response)
}

func (suite *URLGenerationServiceTestSuite) TestGenerateWhenHashExistsShouldReturnExistingShortenedRepSuccessfully() {
	url := "https://www.google.com?q=infracloud"
	urlHash := "559ae209f23cf3866edd0b6bf6c11c8baa128cfc620736d105d582ff8082d1baca509d100daafc165e1ff5d81f1bd3535c9cfa26dba96569a9bec318aee6e756"
	shortenedRep := "FMJRwWE"
	hashMap := models.HashMap{
		HashKey:      urlHash,
		ShortenedRep: shortenedRep,
	}
	urlMap := models.URLMap{
		ShortenedRep: shortenedRep,
		OriginalURL:  url,
	}
	suite.sha512Processor.On("Encode", suite.ctx, url).Return(urlHash)
	suite.hashStoreRepo.On("Fetch", suite.ctx, urlHash).Return(hashMap, nil)
	suite.urlStoreRepo.On("Fetch", suite.ctx, shortenedRep).Return(urlMap, nil)

	response, err := suite.service.Generate(suite.ctx, contracts.GenerateRequest{OriginalURL: url})
	suite.NoError(err)
	suite.Equal(contracts.GenerateResponse{TinyURL: "http://localhost:8181/FMJRwWE"}, response)
}

func TestURLGenerationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(URLGenerationServiceTestSuite))
}
