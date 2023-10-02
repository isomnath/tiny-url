package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type AnalyticsConfigsTestSuite struct {
	suite.Suite
	config *AnalyticsConfigs
}

func (suite *AnalyticsConfigsTestSuite) SetupTest() {
	_ = os.Setenv(topTransformationsPageSize, fmt.Sprintf("%d", 3))
	_ = os.Setenv(topTrafficPageSize, fmt.Sprintf("%d", 4))

	viper.New()
	viper.AutomaticEnv()
	suite.config = analyticsConfigs()
}

func (suite *AnalyticsConfigsTestSuite) TearDownTest() {
	_ = os.Unsetenv(topTransformationsPageSize)
	_ = os.Unsetenv(topTrafficPageSize)
}

func (suite *AnalyticsConfigsTestSuite) TestAllConfigs() {
	suite.Equal(3, suite.config.TopTransformationsPageSize())
	suite.Equal(4, suite.config.TopTrafficPageSize())
}

func TestAnalyticsConfigsTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyticsConfigsTestSuite))
}
