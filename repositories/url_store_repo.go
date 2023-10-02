package repositories

import (
	"context"
	"time"

	"github.com/isomnath/tiny-url/errs"
	"github.com/isomnath/tiny-url/log"
	"github.com/isomnath/tiny-url/models"

	"gopkg.in/redis.v5"
)

type URLStoreRepo struct {
	client *redis.Client
	ttl    time.Duration
}

func (repo *URLStoreRepo) InsertMap(ctx context.Context, model models.URLMap) error {
	err := repo.client.Set(model.ShortenedRep, model.OriginalURL, repo.ttl).Err()
	if err != nil {
		log.Log.WithField("shortened_rep", model.ShortenedRep).
			WithField("original_url", model.OriginalURL).
			WithError(err).Errorf(errs.ErrURLInsertFailure.Error())
		return errs.ErrURLInsertFailure
	}
	return nil
}

func (repo *URLStoreRepo) Fetch(ctx context.Context, shortenedRep string) (models.URLMap, error) {
	statusCmd := repo.client.Get(shortenedRep)
	if statusCmd.Err() != nil {
		if statusCmd.Err() == redis.Nil {
			log.Log.WithField("shortened_rep", shortenedRep).
				WithError(statusCmd.Err()).Errorf(errs.ErrURLFetchFailureKeyDoesNotExist.Error())
			return models.URLMap{}, errs.ErrURLFetchFailureKeyDoesNotExist
		}
		log.Log.WithField("shortened_rep", shortenedRep).
			WithError(statusCmd.Err()).Errorf(errs.ErrURLFetchFailure.Error())
		return models.URLMap{}, errs.ErrURLFetchFailure
	}

	return models.URLMap{ShortenedRep: shortenedRep, OriginalURL: statusCmd.Val()}, nil
}

func NewURLStoreRepo(client *redis.Client, ttl time.Duration) *URLStoreRepo {
	return &URLStoreRepo{
		ttl:    ttl,
		client: client,
	}
}
