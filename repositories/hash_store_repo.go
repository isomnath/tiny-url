package repositories

import (
	"context"
	"time"

	"gopkg.in/redis.v5"

	"github.com/isomnath/tiny-url/errs"
	"github.com/isomnath/tiny-url/log"
	"github.com/isomnath/tiny-url/models"
)

type HashStoreRepo struct {
	client *redis.Client
	ttl    time.Duration
}

func (repo *HashStoreRepo) InsertMap(ctx context.Context, model models.HashMap) error {
	err := repo.client.Set(model.HashKey, model.ShortenedRep, repo.ttl).Err()
	if err != nil {
		log.Log.WithField("hash_key", model.HashKey).
			WithField("shortened_rep", model.ShortenedRep).
			WithError(err).Errorf(errs.ErrHashInsertFailure.Error())
		return errs.ErrHashInsertFailure
	}
	return nil
}

func (repo *HashStoreRepo) Fetch(ctx context.Context, hashKey string) (models.HashMap, error) {
	statusCmd := repo.client.Get(hashKey)
	if statusCmd.Err() != nil {
		if statusCmd.Err() == redis.Nil {
			log.Log.WithField("hash_key", hashKey).
				WithError(statusCmd.Err()).Errorf(errs.ErrHashFetchFailureKeyDoesNotExist.Error())
			return models.HashMap{}, errs.ErrHashFetchFailureKeyDoesNotExist
		}
		log.Log.WithField("hash_key", hashKey).
			WithError(statusCmd.Err()).Errorf(errs.ErrHashFetchFailure.Error())
		return models.HashMap{}, errs.ErrHashFetchFailure
	}

	return models.HashMap{HashKey: hashKey, ShortenedRep: statusCmd.Val()}, nil
}

func NewHashStoreRepo(client *redis.Client, ttl time.Duration) *HashStoreRepo {
	return &HashStoreRepo{
		client: client,
		ttl:    ttl,
	}
}
