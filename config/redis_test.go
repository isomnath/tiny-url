package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type RedisConfigTestSuite struct {
	suite.Suite
	config *RedisConfig
}

func (suite *RedisConfigTestSuite) SetupTest() {
	_ = os.Setenv(redisHost, fmt.Sprintf("%s", "localhost"))
	_ = os.Setenv(redisPort, fmt.Sprintf("%s", "6379"))

	viper.New()
	viper.AutomaticEnv()
	suite.config = redisConfig()
}

func (suite *RedisConfigTestSuite) TearDownTest() {
	_ = os.Unsetenv(redisHost)
	_ = os.Unsetenv(redisPort)
}

func (suite *RedisConfigTestSuite) TestAllConfigs() {
	suite.Equal("localhost", suite.config.Host())
	suite.Equal(6379, suite.config.Port())
}

func TestRedisConfigTestSuite(t *testing.T) {
	suite.Run(t, new(RedisConfigTestSuite))
}
