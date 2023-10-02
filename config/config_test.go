package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	config *Config
}

func (suite *ConfigTestSuite) TestBaseConfigs() {
	LoadBaseConfig()

	suite.Equal("tiny-url", GetAppName())
	suite.Equal("0.0.1", GetAppVersion())
	suite.Equal("staging", GetAppEnvironment())
	suite.Equal(8181, GetAppWebPort())
	suite.Equal("http://localhost:8181", GetAppDNS())
	suite.Equal("DEBUG", GetAppLogLevel())
	suite.Equal("/ping", GetAppHealthCheckAPIPath())
}

func (suite *ConfigTestSuite) TestAllConfigs() {
	LoadBaseConfig()
	LoadRedisConfig()
	LoadAnalyticsConfig()

	suite.Equal("tiny-url", GetAppName())
	suite.Equal("0.0.1", GetAppVersion())
	suite.Equal("staging", GetAppEnvironment())
	suite.Equal(8181, GetAppWebPort())
	suite.Equal("http://localhost:8181", GetAppDNS())
	suite.Equal("DEBUG", GetAppLogLevel())
	suite.Equal("/ping", GetAppHealthCheckAPIPath())

	suite.Equal("localhost", GetRedisConfig().Host())
	suite.Equal(6379, GetRedisConfig().Port())

	suite.Equal(3, GetAnalyticsConfigs().TopTransformationsPageSize())
	suite.Equal(3, GetAnalyticsConfigs().TopTrafficPageSize())
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
