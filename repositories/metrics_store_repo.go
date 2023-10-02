package repositories

import (
	"context"

	"github.com/isomnath/tiny-url/errs"
	"github.com/isomnath/tiny-url/log"
	"github.com/isomnath/tiny-url/models"

	"gopkg.in/redis.v5"
)

const (
	domainsTransformed = "domains_transformed"
	domainsRedirection = "domains_redirection"
)

type MetricsStoreRepo struct {
	client *redis.Client
}

func (repo *MetricsStoreRepo) IncrementTransformationCounter(ctx context.Context, domain string) error {
	err := repo.client.ZIncr(domainsTransformed, redis.Z{Score: 1, Member: domain}).Err()
	if err != nil {
		log.Log.WithField("domain", domain).
			WithError(err).Errorf(errs.ErrMetricsDTIncrFailure.Error())
		return errs.ErrMetricsDTIncrFailure
	}
	return nil
}

func (repo *MetricsStoreRepo) FetchHighestTransformationSet(ctx context.Context, pageSize int64) ([]models.DomainTransformationCounter, error) {
	statusCmd := repo.client.ZRevRangeWithScores(domainsTransformed, 0, pageSize-1)
	if statusCmd.Err() != nil {
		log.Log.WithError(statusCmd.Err()).Errorf(errs.ErrMetricsDTFetchFailure.Error())
		return []models.DomainTransformationCounter{}, errs.ErrMetricsDTFetchFailure
	}

	var page []models.DomainTransformationCounter
	for _, res := range statusCmd.Val() {
		page = append(page, models.DomainTransformationCounter{
			Domain:          res.Member.(string),
			Transformations: res.Score,
		})
	}
	return page, nil
}

func (repo *MetricsStoreRepo) IncrementRedirectionTrafficCounter(ctx context.Context, domain string) error {
	err := repo.client.ZIncr(domainsRedirection, redis.Z{Score: 1, Member: domain}).Err()
	if err != nil {
		log.Log.WithField("domain", domain).
			WithError(err).Errorf(errs.ErrMetricsDRIncrFailure.Error())
		return errs.ErrMetricsDRIncrFailure
	}
	return nil
}

func (repo *MetricsStoreRepo) FetchHighestTrafficSet(ctx context.Context, pageSize int64) ([]models.DomainRedirectionCounter, error) {
	statusCmd := repo.client.ZRevRangeWithScores(domainsRedirection, 0, pageSize-1)
	if statusCmd.Err() != nil {
		log.Log.WithError(statusCmd.Err()).Errorf(errs.ErrMetricsDRFetchFailure.Error())
		return []models.DomainRedirectionCounter{}, errs.ErrMetricsDRFetchFailure
	}

	var page []models.DomainRedirectionCounter
	for _, res := range statusCmd.Val() {
		page = append(page, models.DomainRedirectionCounter{
			Domain:  res.Member.(string),
			Traffic: res.Score,
		})
	}
	return page, nil
}

func NewMetricsStoreRepo(client *redis.Client) *MetricsStoreRepo {
	return &MetricsStoreRepo{
		client: client,
	}
}
