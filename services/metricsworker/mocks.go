package metricsworker

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type metricsStoreRepoMock struct {
	mock.Mock
}

func (mock *metricsStoreRepoMock) IncrementTransformationCounter(ctx context.Context, domain string) error {
	args := mock.Called(ctx, domain)
	return args.Error(0)
}

func (mock *metricsStoreRepoMock) IncrementRedirectionTrafficCounter(ctx context.Context, domain string) error {
	args := mock.Called(ctx, domain)
	return args.Error(0)
}
