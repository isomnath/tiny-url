package metricsworker

import (
	"context"
)

type metricsStoreRepo interface {
	IncrementTransformationCounter(ctx context.Context, domain string) error
	IncrementRedirectionTrafficCounter(ctx context.Context, domain string) error
}
