package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	name               string
	version            string
	environment        string
	webPort            int
	dns                string
	logLevel           string
	healthCheckAPIPath string
	redis              *RedisConfig
	analytics          *AnalyticsConfigs
}

var baseConfig *Config

func LoadBaseConfig() {
	viper.AutomaticEnv()

	viper.SetConfigName("application")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../../")
	viper.AddConfigPath("../../../../")
	viper.AddConfigPath("../../../../../")
	viper.SetConfigType("yaml")
	_ = viper.ReadInConfig()

	baseConfig = &Config{
		name:               getString(appName, true),
		version:            getString(appVersion, true),
		environment:        getString(appEnvironment, true),
		webPort:            getInt(appWebPort, false),
		dns:                getString(appDNS, true),
		logLevel:           getString(appLogLevel, true),
		healthCheckAPIPath: getString(appHealthCheckAPIPath, false),
	}
}

func LoadRedisConfig() {
	baseConfig.redis = redisConfig()
}

func LoadAnalyticsConfig() {
	baseConfig.analytics = analyticsConfigs()
}

func GetAppName() string {
	return baseConfig.name
}

func GetAppVersion() string {
	return baseConfig.version
}

func GetAppEnvironment() string {
	return baseConfig.environment
}

func GetAppWebPort() int {
	return baseConfig.webPort
}

func GetAppDNS() string {
	return baseConfig.dns
}

func GetAppLogLevel() string {
	return baseConfig.logLevel
}

func GetAppHealthCheckAPIPath() string {
	return baseConfig.healthCheckAPIPath
}

func GetRedisConfig() *RedisConfig {
	return baseConfig.redis
}

func GetAnalyticsConfigs() *AnalyticsConfigs {
	return baseConfig.analytics
}
