package dependencies

import (
	"fmt"

	"gopkg.in/redis.v5"

	"github.com/isomnath/tiny-url/config"
)

type Dependencies struct {
	Controllers *Controllers
}

func InitializeDependencies() *Dependencies {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.GetRedisConfig().Host(), config.GetRedisConfig().Port()),
		DB:   0,
	})

	client.FlushAll()

	return &Dependencies{
		Controllers: initializeControllers(initializeServices(initializeRepositories(client))),
	}
}
