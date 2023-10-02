package repositories

import (
	"context"
	"math/rand"
	"time"

	"gopkg.in/redis.v5"

	"github.com/isomnath/tiny-url/errs"
	"github.com/isomnath/tiny-url/log"
	"github.com/isomnath/tiny-url/models"
)

const atomicCounter = "atomic_counter"

type SeriesCounterRepo struct {
	client *redis.Client
	ttl    time.Duration
}

func (repo *SeriesCounterRepo) Init(ctx context.Context) error {
	num := rand.Intn(1234567890)
	sCmd := repo.client.Set(atomicCounter, num, repo.ttl)
	if sCmd.Err() != nil {
		log.Log.WithError(sCmd.Err()).Errorf(errs.ErrSeriesCountInitFailure.Error())
		return errs.ErrSeriesCountInitFailure
	}
	return nil
}

func (repo *SeriesCounterRepo) Next(ctx context.Context) (models.Series, error) {
	iCmd := repo.client.Incr(atomicCounter)
	if iCmd.Err() != nil {
		log.Log.WithError(iCmd.Err()).Error(errs.ErrSeriesCounterFailure.Error())
		return models.Series(0), errs.ErrSeriesCounterFailure
	}
	series := iCmd.Val()
	return models.Series(series), nil
}

func NewSeriesCounterRepo(client *redis.Client, ttl time.Duration) *SeriesCounterRepo {
	return &SeriesCounterRepo{
		client: client,
		ttl:    ttl,
	}
}
