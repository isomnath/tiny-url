package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/log"
	"github.com/isomnath/tiny-url/models"

	"github.com/stretchr/testify/suite"
	"gopkg.in/redis.v5"
)

type MetricsStoreRepoTestSuite struct {
	suite.Suite
	ctx    context.Context
	client *redis.Client
	repo   *MetricsStoreRepo
}

func (suite *MetricsStoreRepoTestSuite) SetupTest() {
	suite.ctx = context.Background()
	config.LoadBaseConfig()
	config.LoadRedisConfig()
	log.Setup()
	suite.client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.GetRedisConfig().Host(), config.GetRedisConfig().Port()),
		DB:   0,
	})
	suite.repo = NewMetricsStoreRepo(suite.client)
}

func (suite *MetricsStoreRepoTestSuite) TearDownTest() {
	suite.client.FlushDb()
}

func (suite *MetricsStoreRepoTestSuite) TestIncrementTransformationCounterReturnError() {
	_ = suite.client.Close()
	err := suite.repo.IncrementTransformationCounter(suite.ctx, "Google")
	suite.Equal("failed to increment domain transformation counter", err.Error())
}

func (suite *MetricsStoreRepoTestSuite) TestIncrementTransformationCounterSuccess() {
	err := suite.repo.IncrementTransformationCounter(suite.ctx, "Google")
	suite.NoError(err)
}

func (suite *MetricsStoreRepoTestSuite) TestFetchHighestTransformationSetReturnError() {
	_ = suite.client.Close()
	_, err := suite.repo.FetchHighestTransformationSet(suite.ctx, 3)
	suite.Equal("failed to fetch list of most transformed domains", err.Error())
}

func (suite *MetricsStoreRepoTestSuite) TestFetchHighestTransformationSetSuccess() {
	_ = suite.client.ZIncrBy(domainsTransformed, 2, "Google")
	_ = suite.client.ZIncrBy(domainsTransformed, 5, "Youtube")
	_ = suite.client.ZIncrBy(domainsTransformed, 4, "Stackoverflow")
	_ = suite.client.ZIncrBy(domainsTransformed, 7, "Substack")

	expectedCounter := []models.DomainTransformationCounter{
		{
			Domain:          "Substack",
			Transformations: 7,
		},
		{
			Domain:          "Youtube",
			Transformations: 5,
		},
		{
			Domain:          "Stackoverflow",
			Transformations: 4,
		},
	}
	counter, err := suite.repo.FetchHighestTransformationSet(suite.ctx, 3)
	suite.NoError(err)
	suite.Equal(expectedCounter, counter)
}

func (suite *MetricsStoreRepoTestSuite) TestIncrementRedirectionTrafficCounterReturnError() {
	_ = suite.client.Close()
	err := suite.repo.IncrementRedirectionTrafficCounter(suite.ctx, "Google")
	suite.Equal("failed to increment domain traffic counter", err.Error())
}

func (suite *MetricsStoreRepoTestSuite) TestIncrementRedirectionTrafficCounterSuccess() {
	err := suite.repo.IncrementRedirectionTrafficCounter(suite.ctx, "Google")
	suite.NoError(err)
}

func (suite *MetricsStoreRepoTestSuite) TestFetchHighestTrafficSetReturnError() {
	_ = suite.client.Close()
	_, err := suite.repo.FetchHighestTrafficSet(suite.ctx, 3)
	suite.Equal("failed to fetch list of domains with most traffic", err.Error())
}

func (suite *MetricsStoreRepoTestSuite) TestFetchHighestTrafficSetSuccess() {
	_ = suite.client.ZIncrBy(domainsRedirection, 2, "Google")
	_ = suite.client.ZIncrBy(domainsRedirection, 5, "Youtube")
	_ = suite.client.ZIncrBy(domainsRedirection, 4, "Stackoverflow")
	_ = suite.client.ZIncrBy(domainsRedirection, 7, "Substack")

	expectedCounter := []models.DomainRedirectionCounter{
		{
			Domain:  "Substack",
			Traffic: 7,
		},
		{
			Domain:  "Youtube",
			Traffic: 5,
		},
		{
			Domain:  "Stackoverflow",
			Traffic: 4,
		},
	}
	counter, err := suite.repo.FetchHighestTrafficSet(suite.ctx, 3)
	suite.NoError(err)
	suite.Equal(expectedCounter, counter)
}

func TestMetricsStoreRepoTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsStoreRepoTestSuite))
}
