package controllers

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/isomnath/tiny-url/contracts"
)

type urlGenerationServiceMock struct {
	mock.Mock
}

func (mock *urlGenerationServiceMock) Generate(ctx context.Context, request contracts.GenerateRequest) (contracts.GenerateResponse, error) {
	args := mock.Called(request)
	return args.Get(0).(contracts.GenerateResponse), args.Error(1)
}

type urlRedirectionServiceMock struct {
	mock.Mock
}

func (mock *urlRedirectionServiceMock) Redirect(ctx context.Context, shortenedRep string) (string, error) {
	args := mock.Called(shortenedRep)
	return args.String(0), args.Error(1)
}

type analyticsServiceMock struct {
	mock.Mock
}

func (mock *analyticsServiceMock) FetchTopDomainsByTransformations(ctx context.Context) (contracts.TransformationResponse, error) {
	args := mock.Called()
	return args.Get(0).(contracts.TransformationResponse), args.Error(1)
}

func (mock *analyticsServiceMock) FetchTopDomainsByTraffic(ctx context.Context) (contracts.TrafficResponse, error) {
	args := mock.Called()
	return args.Get(0).(contracts.TrafficResponse), args.Error(1)
}
