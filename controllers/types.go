package controllers

import (
	"context"

	"github.com/isomnath/tiny-url/contracts"
)

type urlGenerationService interface {
	Generate(ctx context.Context, request contracts.GenerateRequest) (contracts.GenerateResponse, error)
}

type urlRedirectionService interface {
	Redirect(ctx context.Context, shortenedRep string) (string, error)
}

type analyticsService interface {
	FetchTopDomainsByTransformations(ctx context.Context) (contracts.TransformationResponse, error)
	FetchTopDomainsByTraffic(ctx context.Context) (contracts.TrafficResponse, error)
}
