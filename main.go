package main

import (
	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/dependencies"
	"github.com/isomnath/tiny-url/log"
	"github.com/isomnath/tiny-url/server"
)

func main() {
	config.LoadBaseConfig()
	config.LoadRedisConfig()
	config.LoadAnalyticsConfig()
	log.Setup()
	dep := dependencies.InitializeDependencies()
	router := server.InitializeRouter(dep)
	server.StartServer(router)
}
