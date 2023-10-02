package controllers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/errs"
	"github.com/isomnath/tiny-url/log"
)

type URLRedirectionControllerTestSuite struct {
	suite.Suite
	router                *mux.Router
	urlRedirectionService *urlRedirectionServiceMock
	ctrl                  *URLRedirectionController
}

func (suite *URLRedirectionControllerTestSuite) SetupTest() {
	config.LoadBaseConfig()
	log.Setup()
	suite.urlRedirectionService = new(urlRedirectionServiceMock)
	suite.ctrl = NewURLRedirectionController(suite.urlRedirectionService)
	suite.router = mux.NewRouter()
}

func (suite *URLRedirectionControllerTestSuite) TestRedirectReturnServerError() {
	suite.router.HandleFunc("/{tiny-url}", suite.ctrl.Redirect).Methods(http.MethodGet)

	suite.urlRedirectionService.On("Redirect", "01LY7VK").Return("", errs.ErrFailedToFetchURLFromURLStore)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/01LY7VK", nil)
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":false,"error":{"message":"failed to fetch shortened url from url store"}}`
	suite.Equal(http.StatusInternalServerError, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func (suite *URLRedirectionControllerTestSuite) TestRedirectSuccessfullyRedirectsToOriginalURL() {
	suite.router.HandleFunc("/{tiny-url}", suite.ctrl.Redirect).Methods(http.MethodGet)
	response := "https://www.google.com?q=infracloud"

	suite.urlRedirectionService.On("Redirect", "01LY7VK").Return(response, nil)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/01LY7VK", nil)
	suite.router.ServeHTTP(recorder, request)

	redirectedURL := recorder.Result().Header.Get("Location")
	suite.Equal(http.StatusOK, recorder.Code)
	suite.Equal(response, redirectedURL)
}

func TestURLRedirectionControllerTestSuite(t *testing.T) {
	suite.Run(t, new(URLRedirectionControllerTestSuite))
}
