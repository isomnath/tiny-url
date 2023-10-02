package metricsworker

import (
	"context"
	"errors"
	"testing"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/log"

	"github.com/stretchr/testify/suite"
)

type MetricsWorkerTestSuite struct {
	suite.Suite
	ctx                  context.Context
	metricsStoreRepoMock *metricsStoreRepoMock
	worker               *MetricsWorker
}

func (suite *MetricsWorkerTestSuite) SetupTest() {
	config.LoadBaseConfig()
	log.Setup()
	suite.ctx = context.Background()
	suite.metricsStoreRepoMock = new(metricsStoreRepoMock)
	suite.worker = NewMetricsWorker(suite.metricsStoreRepoMock)
}

func (suite *MetricsWorkerTestSuite) TestIncrementGenerationCounterError() {
	suite.metricsStoreRepoMock.On("IncrementTransformationCounter", suite.ctx, "Google").Return(errors.New("redis error"))
	suite.worker.IncrementGenerationCounter(suite.ctx, "Google")
}

func (suite *MetricsWorkerTestSuite) TestIncrementGenerationCounterSuccess() {
	suite.metricsStoreRepoMock.On("IncrementTransformationCounter", suite.ctx, "Google").Return(nil)
	suite.worker.IncrementGenerationCounter(suite.ctx, "Google")
}

func (suite *MetricsWorkerTestSuite) TestIncrementRedirectionTrafficCounterError() {
	suite.metricsStoreRepoMock.On("IncrementRedirectionTrafficCounter", suite.ctx, "Google").Return(errors.New("redis error"))
	suite.worker.IncrementRedirectionTrafficCounter(suite.ctx, "Google")
}

func (suite *MetricsWorkerTestSuite) TestIncrementRedirectionTrafficCounterSuccess() {
	suite.metricsStoreRepoMock.On("IncrementRedirectionTrafficCounter", suite.ctx, "Google").Return(nil)
	suite.worker.IncrementRedirectionTrafficCounter(suite.ctx, "Google")
}

func TestMetricsWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsWorkerTestSuite))
}
