package dependencies

import (
	"context"
	"time"

	"gopkg.in/redis.v5"

	"github.com/isomnath/tiny-url/log"
	"github.com/isomnath/tiny-url/repositories"
)

type Repositories struct {
	HashStoreRepo     *repositories.HashStoreRepo
	URLStoreRepo      *repositories.URLStoreRepo
	MetricsStoreRepo  *repositories.MetricsStoreRepo
	SeriesCounterRepo *repositories.SeriesCounterRepo
}

func initializeRepositories(client *redis.Client) *Repositories {
	seriesCounterRepo := repositories.NewSeriesCounterRepo(client, 3650*24*time.Hour)
	err := seriesCounterRepo.Init(context.Background())
	if err != nil {
		log.Log.WithError(err).Fatalf(err.Error())
	}
	return &Repositories{
		HashStoreRepo:     repositories.NewHashStoreRepo(client, 3650*24*time.Hour),
		URLStoreRepo:      repositories.NewURLStoreRepo(client, 3650*24*time.Hour),
		MetricsStoreRepo:  repositories.NewMetricsStoreRepo(client),
		SeriesCounterRepo: seriesCounterRepo,
	}
}
