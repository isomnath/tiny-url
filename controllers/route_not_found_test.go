package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/log"

	"github.com/stretchr/testify/suite"
)

type RouteNotFoundHandlerTestSuite struct {
	suite.Suite
}

func (suite *RouteNotFoundHandlerTestSuite) SetupSuite() {
	config.LoadBaseConfig()
	log.Setup()
}

func (suite *RouteNotFoundHandlerTestSuite) TestRouteNotFoundHandlerShouldReturnNotFound() {
	rw := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/invalid_path", nil)
	suite.NoError(err, "failed to create a request")

	RouteNotFoundHandler(rw, r)

	suite.Equal(http.StatusNotFound, rw.Code)
	suite.Equal("{\"success\":false,\"error\":{\"message\":\"route /invalid_path not found\"}}", rw.Body.String())
}

func TestRouteNotFoundHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(RouteNotFoundHandlerTestSuite))
}
