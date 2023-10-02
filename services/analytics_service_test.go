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

type AnalyticsServiceTestSuite struct {
	suite.Suite
	ctx                  context.Context
	metricsStoreRepoMock *metricsStoreRepoMock
	service              *AnalyticsService
}

func (suite *AnalyticsServiceTestSuite) SetupTest() {
	config.LoadBaseConfig()
	config.LoadAnalyticsConfig()
	log.Setup()
	suite.ctx = context.Background()
	suite.metricsStoreRepoMock = new(metricsStoreRepoMock)
	suite.service = NewAnalyticsService(suite.metricsStoreRepoMock)
}

func (suite *AnalyticsServiceTestSuite) TestFetchTopDomainsByTransformationsReturnError() {
	suite.metricsStoreRepoMock.On("FetchHighestTransformationSet", suite.ctx,
		int64(config.GetAnalyticsConfigs().TopTransformationsPageSize())).Return([]models.DomainTransformationCounter{}, errs.ErrMetricsDTFetchFailure)

	response, err := suite.service.FetchTopDomainsByTransformations(suite.ctx)
	suite.Equal(errs.ErrMetricsDTFetchFailure, err)
	suite.Empty(response)
}

func (suite *AnalyticsServiceTestSuite) TestFetchTopDomainsByTransformationsReturnSuccessResponse() {
	transformationSet := []models.DomainTransformationCounter{
		{
			Domain:          "google",
			Transformations: 4,
		},
		{
			Domain:          "youtube",
			Transformations: 3,
		},
		{
			Domain:          "stackoverflow",
			Transformations: 2,
		},
	}
	expectedResponse := contracts.TransformationResponse{
		TopTransformations: []contracts.Transformation{
			{
				Domain: "google",
				Count:  4,
			},
			{
				Domain: "youtube",
				Count:  3,
			},
			{
				Domain: "stackoverflow",
				Count:  2,
			},
		},
	}

	suite.metricsStoreRepoMock.On("FetchHighestTransformationSet", suite.ctx,
		int64(config.GetAnalyticsConfigs().TopTransformationsPageSize())).Return(transformationSet, nil)

	response, err := suite.service.FetchTopDomainsByTransformations(suite.ctx)
	suite.NoError(err)
	suite.Equal(expectedResponse, response)
}

func (suite *AnalyticsServiceTestSuite) TestFetchTopDomainsByTrafficReturnError() {
	suite.metricsStoreRepoMock.On("FetchHighestTrafficSet", suite.ctx,
		int64(config.GetAnalyticsConfigs().TopTrafficPageSize())).Return([]models.DomainRedirectionCounter{}, errs.ErrMetricsDRFetchFailure)
	response, err := suite.service.FetchTopDomainsByTraffic(suite.ctx)
	suite.Equal(errs.ErrMetricsDRFetchFailure, err)
	suite.Empty(response)
}

func (suite *AnalyticsServiceTestSuite) TestFetchTopDomainsByTrafficReturnSuccessResponse() {
	trafficSet := []models.DomainRedirectionCounter{
		{
			Domain:  "google",
			Traffic: 4,
		},
		{
			Domain:  "youtube",
			Traffic: 3,
		},
		{
			Domain:  "stackoverflow",
			Traffic: 2,
		},
	}
	expectedResponse := contracts.TrafficResponse{
		TopTraffic: []contracts.Traffic{
			{
				Domain: "google",
				Count:  4,
			},
			{
				Domain: "youtube",
				Count:  3,
			},
			{
				Domain: "stackoverflow",
				Count:  2,
			},
		},
	}

	suite.metricsStoreRepoMock.On("FetchHighestTrafficSet", suite.ctx,
		int64(config.GetAnalyticsConfigs().TopTrafficPageSize())).Return(trafficSet, nil)
	response, err := suite.service.FetchTopDomainsByTraffic(suite.ctx)
	suite.NoError(err)
	suite.Equal(expectedResponse, response)
}

func TestAnalyticsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyticsServiceTestSuite))
}
