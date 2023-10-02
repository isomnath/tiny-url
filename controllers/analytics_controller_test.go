package controllers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/isomnath/tiny-url/contracts"
	"github.com/isomnath/tiny-url/errs"
)

type AnalyticsControllerTestSuite struct {
	suite.Suite
	router           *mux.Router
	analyticsService *analyticsServiceMock
	ctrl             *AnalyticsController
}

func (suite *AnalyticsControllerTestSuite) SetupTest() {
	suite.router = mux.NewRouter()
	suite.analyticsService = new(analyticsServiceMock)
	suite.ctrl = NewAnalyticsController(suite.analyticsService)
}

func (suite *AnalyticsControllerTestSuite) TestTopHighTransformedDomainsReturnServerError() {
	suite.router.HandleFunc("/v1/analytics/domains/highest_transformation", suite.ctrl.TopHighTransformedDomains).Methods(http.MethodGet)

	suite.analyticsService.On("FetchTopDomainsByTransformations").Return(contracts.TransformationResponse{}, errs.ErrMetricsDTFetchFailure)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/analytics/domains/highest_transformation", nil)
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":false,"error":{"message":"failed to fetch list of most transformed domains"}}`
	suite.Equal(http.StatusInternalServerError, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func (suite *AnalyticsControllerTestSuite) TestTopHighTransformedDomainsReturnSuccessResponse() {
	suite.router.HandleFunc("/v1/analytics/domains/highest_transformation", suite.ctrl.TopHighTransformedDomains).Methods(http.MethodGet)

	response := contracts.TransformationResponse{
		TopTransformations: []contracts.Transformation{
			{
				Domain: "google",
				Count:  4,
			},
			{
				Domain: "youtube",
				Count:  3,
			},
			{
				Domain: "stackoverflow",
				Count:  2,
			},
		},
	}
	suite.analyticsService.On("FetchTopDomainsByTransformations").Return(response, nil)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/analytics/domains/highest_transformation", nil)
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":true,"data":{"top_transformations":[{"domain":"google","count":4},{"domain":"youtube","count":3},{"domain":"stackoverflow","count":2}]}}`
	suite.Equal(http.StatusOK, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func (suite *AnalyticsControllerTestSuite) TestTopHighTrafficDomainsReturnServerError() {
	suite.router.HandleFunc("/v1/analytics/domains/highest_traffic", suite.ctrl.TopHighTrafficDomains).Methods(http.MethodGet)

	suite.analyticsService.On("FetchTopDomainsByTraffic").Return(contracts.TrafficResponse{}, errs.ErrMetricsDRFetchFailure)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/analytics/domains/highest_traffic", nil)
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":false,"error":{"message":"failed to fetch list of domains with most traffic"}}`
	suite.Equal(http.StatusInternalServerError, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func (suite *AnalyticsControllerTestSuite) TestTopHighTrafficDomainsReturnSuccessResponse() {
	suite.router.HandleFunc("/v1/analytics/domains/highest_traffic", suite.ctrl.TopHighTrafficDomains).Methods(http.MethodGet)

	response := contracts.TrafficResponse{
		TopTraffic: []contracts.Traffic{
			{
				Domain: "google",
				Count:  4,
			},
			{
				Domain: "youtube",
				Count:  3,
			},
			{
				Domain: "stackoverflow",
				Count:  2,
			},
		},
	}
	suite.analyticsService.On("FetchTopDomainsByTraffic").Return(response, nil)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/analytics/domains/highest_traffic", nil)
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":true,"data":{"top_traffic":[{"domain":"google","count":4},{"domain":"youtube","count":3},{"domain":"stackoverflow","count":2}]}}`
	suite.Equal(http.StatusOK, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func TestAnalyticsControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyticsControllerTestSuite))
}
