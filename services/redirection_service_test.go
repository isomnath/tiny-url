package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/isomnath/tiny-url/errs"
	"github.com/isomnath/tiny-url/models"
)

type URLRedirectionServiceTestSuite struct {
	suite.Suite
	ctx               context.Context
	urlStoreRepoMock  *urlStoreRepoMock
	urlProcessorMock  *urlProcessorMock
	metricsWorkerMock *metricsWorkerMock
	service           *URLRedirectionService
}

func (suite *URLRedirectionServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.urlStoreRepoMock = new(urlStoreRepoMock)
	suite.urlProcessorMock = new(urlProcessorMock)
	suite.metricsWorkerMock = new(metricsWorkerMock)
	suite.service = NewURLRedirectionService(suite.urlStoreRepoMock, suite.urlProcessorMock, suite.metricsWorkerMock)
}

func (suite *URLRedirectionServiceTestSuite) TestRedirectReturnErrorWhenShortenedURLNotFound() {
	shortenedRep := "FMJRwWE"
	suite.urlStoreRepoMock.On("Fetch", suite.ctx, shortenedRep).Return(models.URLMap{}, errs.ErrURLFetchFailureKeyDoesNotExist)
	originalURL, err := suite.service.Redirect(suite.ctx, shortenedRep)
	suite.Equal(errs.ErrURLDoesNotExist, err)
	suite.Equal("", originalURL)
}

func (suite *URLRedirectionServiceTestSuite) TestRedirectReturnErrorWhenFailedToFetchURLFromStore() {
	shortenedRep := "FMJRwWE"
	suite.urlStoreRepoMock.On("Fetch", suite.ctx, shortenedRep).Return(models.URLMap{}, errs.ErrURLFetchFailure)
	originalURL, err := suite.service.Redirect(suite.ctx, shortenedRep)
	suite.Equal(errs.ErrFailedToFetchURLFromURLStore, err)
	suite.Equal("", originalURL)
}

func (suite *URLRedirectionServiceTestSuite) TestRedirectReturnOriginalURLSuccessfully() {
	shortenedRep := "FMJRwWE"
	url := "https://www.google.com?q=infracloud"
	urlMap := models.URLMap{ShortenedRep: shortenedRep, OriginalURL: url}
	suite.urlStoreRepoMock.On("Fetch", suite.ctx, shortenedRep).Return(urlMap, nil)
	suite.urlProcessorMock.On("ExtractDomain", suite.ctx, url).Return("google")
	suite.metricsWorkerMock.On("IncrementRedirectionTrafficCounter", suite.ctx, "google")
	originalURL, err := suite.service.Redirect(suite.ctx, shortenedRep)
	suite.NoError(err)
	suite.Equal(url, originalURL)
}

func TestURLRedirectionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(URLRedirectionServiceTestSuite))
}
