package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/log"
	"github.com/isomnath/tiny-url/models"

	"github.com/stretchr/testify/suite"
	"gopkg.in/redis.v5"
)

type URLStoreRepoTestSuite struct {
	suite.Suite
	ctx    context.Context
	client *redis.Client
	repo   *URLStoreRepo
}

func (suite *URLStoreRepoTestSuite) SetupTest() {
	suite.ctx = context.Background()
	config.LoadBaseConfig()
	config.LoadRedisConfig()
	log.Setup()
	suite.client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.GetRedisConfig().Host(), config.GetRedisConfig().Port()),
		DB:   0,
	})
	suite.repo = NewURLStoreRepo(suite.client, 10*time.Minute)
}

func (suite *URLStoreRepoTestSuite) TearDownTest() {
	suite.client.FlushDb()
}

func (suite *URLStoreRepoTestSuite) TestInsertMapReturnError() {
	_ = suite.client.Close()
	err := suite.repo.InsertMap(suite.ctx, models.URLMap{ShortenedRep: "3wl0boM", OriginalURL: "https://www.youtube.com/some-page-123"})
	suite.Equal("failed to insert url mapping", err.Error())
}

func (suite *URLStoreRepoTestSuite) TestInsertMapSuccess() {
	err := suite.repo.InsertMap(suite.ctx, models.URLMap{ShortenedRep: "3wl0boM", OriginalURL: "https://www.youtube.com/some-page-123"})
	suite.NoError(err)
}

func (suite *URLStoreRepoTestSuite) TestFetchReturnErrorWhenKeyDoesNotExist() {
	mp, err := suite.repo.Fetch(suite.ctx, "3wl0boM")
	suite.Empty(mp)
	suite.Equal("failed to fetch url mapping as shortened url does not exist in the set", err.Error())
}

func (suite *URLStoreRepoTestSuite) TestFetchReturnError() {
	_ = suite.client.Close()
	mp, err := suite.repo.Fetch(suite.ctx, "3wl0boM")
	suite.Empty(mp)
	suite.Equal("failed to fetch url mapping", err.Error())
}

func (suite *URLStoreRepoTestSuite) TestFetchReturnSuccess() {
	_ = suite.repo.InsertMap(suite.ctx, models.URLMap{ShortenedRep: "3wl0boM", OriginalURL: "https://www.youtube.com/some-page-123"})
	mp, err := suite.repo.Fetch(suite.ctx, "3wl0boM")
	suite.NoError(err)
	suite.Equal(models.URLMap{ShortenedRep: "3wl0boM", OriginalURL: "https://www.youtube.com/some-page-123"}, mp)
}

func TestURLStoreRepoTestSuite(t *testing.T) {
	suite.Run(t, new(URLStoreRepoTestSuite))
}
