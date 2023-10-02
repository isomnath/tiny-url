package log

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/isomnath/tiny-url/config"
)

type LoggerTestSuite struct {
	suite.Suite
}

func (suite *LoggerTestSuite) SetupSuite() {
	config.LoadBaseConfig()
}

func (suite *LoggerTestSuite) TestLoggerPanic() {
	_ = os.Setenv("APP_LOG_LEVEL", "INVALID")
	config.LoadBaseConfig()
	suite.Panics(func() {
		Setup()
	})
}

func (suite *LoggerTestSuite) TestLogger() {
	Setup()

	Log.Fatalf("test message, args1: %s", "123")
	Log.Errorf("test message, args1: %s", "123")
	Log.Warnf("test message, args1: %s", "123")
	Log.Infof("test message, args1: %s", "123")

	request, _ := http.NewRequest(http.MethodGet, "/test/path/123", nil)
	startTime := time.Now()
	responseTime := startTime.Add(20 * time.Millisecond)
	Log.HTTPStatInfo(request, startTime, responseTime, http.StatusCreated)
	Log.HTTPErrorf(request, "test message, args1: %s", "123")
	Log.HTTPWarnf(request, "test message, args1: %s", "123")
	Log.HTTPInfof(request, "test message, args1: %s", "123")

	//Log.RedisFatalf("test message, args1: %s", "123")
	Log.RedisErrorf("test message, args1: %s", "123")
	Log.RedisWarnf("test message, args1: %s", "123")
	Log.RedisInfof("test message, args1: %s", "123")
}

func (suite *LoggerTestSuite) TestRuntimeProcCaptureFailure() {
	expectedFields := logrus.Fields{"file": "unknown", "function": "unknown"}
	actualFields := Log.getProcessFields(uintptr(1), "", 0, false)
	suite.Equal(expectedFields, actualFields)
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
