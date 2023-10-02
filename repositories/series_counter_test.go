package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gopkg.in/redis.v5"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/log"
	"github.com/isomnath/tiny-url/models"
)

type SeriesCounterTestSuite struct {
	suite.Suite
	ctx    context.Context
	client *redis.Client
	repo   *SeriesCounterRepo
}

func (suite *SeriesCounterTestSuite) SetupTest() {
	suite.ctx = context.Background()
	config.LoadBaseConfig()
	config.LoadRedisConfig()
	log.Setup()
	suite.client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.GetRedisConfig().Host(), config.GetRedisConfig().Port()),
		DB:   0,
	})
	suite.repo = NewSeriesCounterRepo(suite.client, 10*time.Minute)
}

func (suite *SeriesCounterTestSuite) TearDownTest() {
	suite.client.FlushDb()
}

func (suite *SeriesCounterTestSuite) TestInitReturnError() {
	_ = suite.client.Close()
	err := suite.repo.Init(suite.ctx)
	suite.Equal("failed to set atomic counter with initial value", err.Error())
}

func (suite *SeriesCounterTestSuite) TestInitSuccess() {
	err := suite.repo.Init(suite.ctx)
	suite.NoError(err)
}

func (suite *SeriesCounterTestSuite) TestNextReturnError() {
	_ = suite.client.Close()
	series, err := suite.repo.Next(suite.ctx)
	suite.Equal("atomic counter could not be incremented", err.Error())
	suite.Equal(models.Series(0), series)
}

func (suite *SeriesCounterTestSuite) TestNextSuccess() {
	series, err := suite.repo.Next(suite.ctx)
	suite.NoError(err)
	suite.Equal(models.Series(1), series)
}

func TestSeriesCounterTestSuite(t *testing.T) {
	suite.Run(t, new(SeriesCounterTestSuite))
}
