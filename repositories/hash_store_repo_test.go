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

type HashStoreRepoTestSuite struct {
	suite.Suite
	ctx    context.Context
	client *redis.Client
	repo   *HashStoreRepo
}

func (suite *HashStoreRepoTestSuite) SetupTest() {
	suite.ctx = context.Background()
	config.LoadBaseConfig()
	config.LoadRedisConfig()
	log.Setup()
	suite.client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.GetRedisConfig().Host(), config.GetRedisConfig().Port()),
		DB:   0,
	})
	suite.repo = NewHashStoreRepo(suite.client, 10*time.Minute)
}

func (suite *HashStoreRepoTestSuite) TearDownTest() {
	suite.client.FlushDb()
}

func (suite *HashStoreRepoTestSuite) TestInsertMapReturnError() {
	_ = suite.client.Close()
	err := suite.repo.InsertMap(suite.ctx, models.HashMap{HashKey: "7fc6dccd467d94f1f63601a33fbebb8519b5fe639a7550d733f6e91d7d6b80b5448463c33853f8591a611e2a6dac09ed943299e79d1b6b3b9e9d9cad3d816135", ShortenedRep: "3wl0boM"})
	suite.Equal("failed to insert hash map", err.Error())
}

func (suite *HashStoreRepoTestSuite) TestInsertMapSuccess() {
	err := suite.repo.InsertMap(suite.ctx, models.HashMap{HashKey: "7fc6dccd467d94f1f63601a33fbebb8519b5fe639a7550d733f6e91d7d6b80b5448463c33853f8591a611e2a6dac09ed943299e79d1b6b3b9e9d9cad3d816135", ShortenedRep: "3wl0boM"})
	suite.NoError(err)
}

func (suite *HashStoreRepoTestSuite) TestFetchReturnErrorWhenKeyDoesNotExist() {
	mp, err := suite.repo.Fetch(suite.ctx, "7fc6dccd467d94f1f63601a33fbebb8519b5fe639a7550d733f6e91d7d6b80b5448463c33853f8591a611e2a6dac09ed943299e79d1b6b3b9e9d9cad3d816135")
	suite.Empty(mp)
	suite.Equal("failed to fetch hash mapping as hash does not exist in the set", err.Error())
}

func (suite *HashStoreRepoTestSuite) TestFetchReturnError() {
	_ = suite.client.Close()
	mp, err := suite.repo.Fetch(suite.ctx, "7fc6dccd467d94f1f63601a33fbebb8519b5fe639a7550d733f6e91d7d6b80b5448463c33853f8591a611e2a6dac09ed943299e79d1b6b3b9e9d9cad3d816135")
	suite.Empty(mp)
	suite.Equal("failed to fetch hash map", err.Error())
}

func (suite *HashStoreRepoTestSuite) TestFetchReturnSuccess() {
	_ = suite.repo.InsertMap(suite.ctx, models.HashMap{HashKey: "7fc6dccd467d94f1f63601a33fbebb8519b5fe639a7550d733f6e91d7d6b80b5448463c33853f8591a611e2a6dac09ed943299e79d1b6b3b9e9d9cad3d816135", ShortenedRep: "3wl0boM"})
	mp, err := suite.repo.Fetch(suite.ctx, "7fc6dccd467d94f1f63601a33fbebb8519b5fe639a7550d733f6e91d7d6b80b5448463c33853f8591a611e2a6dac09ed943299e79d1b6b3b9e9d9cad3d816135")
	suite.NoError(err)
	suite.Equal(models.HashMap{HashKey: "7fc6dccd467d94f1f63601a33fbebb8519b5fe639a7550d733f6e91d7d6b80b5448463c33853f8591a611e2a6dac09ed943299e79d1b6b3b9e9d9cad3d816135", ShortenedRep: "3wl0boM"}, mp)
}

func TestHashStoreRepoTestSuite(t *testing.T) {
	suite.Run(t, new(HashStoreRepoTestSuite))
}
