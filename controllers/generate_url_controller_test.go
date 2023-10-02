package controllers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/contracts"
	"github.com/isomnath/tiny-url/errs"
	"github.com/isomnath/tiny-url/log"
)

type GenerateURLControllerTestSuite struct {
	suite.Suite
	router               *mux.Router
	urlGenerationService *urlGenerationServiceMock
	ctrl                 *GenerateURLController
}

func (suite *GenerateURLControllerTestSuite) SetupTest() {
	config.LoadBaseConfig()
	log.Setup()
	suite.urlGenerationService = new(urlGenerationServiceMock)
	suite.ctrl = NewGenerateURLController(suite.urlGenerationService)
	suite.router = mux.NewRouter()
}

func (suite *GenerateURLControllerTestSuite) TestGenerateReturnBadRequestErrorWhenRequestUnmarshalFails() {
	body := `{"original_url": 124151}`

	suite.router.HandleFunc("/v1/generate", suite.ctrl.Generate).Methods(http.MethodPost)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/generate", bytes.NewBuffer([]byte(body)))
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":false,"error":{"message":"failed to deserialize json request body to destination interface"}}`
	suite.Equal(http.StatusBadRequest, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func (suite *GenerateURLControllerTestSuite) TestGenerateReturnBadRequestErrorWhenURLIsEmpty() {
	body := `{"original_url": ""}`

	suite.router.HandleFunc("/v1/generate", suite.ctrl.Generate).Methods(http.MethodPost)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/generate", bytes.NewBuffer([]byte(body)))
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":false,"error":{"message":"empty url"}}`
	suite.Equal(http.StatusBadRequest, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func (suite *GenerateURLControllerTestSuite) TestGenerateReturnBadRequestErrorWhenURLIsInvalid() {
	body := `{"original_url": "https://%%zz"}`

	suite.router.HandleFunc("/v1/generate", suite.ctrl.Generate).Methods(http.MethodPost)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/generate", bytes.NewBuffer([]byte(body)))
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":false,"error":{"message":"invalid url"}}`
	suite.Equal(http.StatusBadRequest, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func (suite *GenerateURLControllerTestSuite) TestGenerateReturnServerError() {
	body := `{"original_url": "https://www.google.com?q=infracloud"}`
	deserializedRequest := contracts.GenerateRequest{OriginalURL: "https://www.google.com?q=infracloud"}

	suite.urlGenerationService.On("Generate", deserializedRequest).Return(contracts.GenerateResponse{}, errs.ErrFailedToSaveURLInStore)
	suite.router.HandleFunc("/v1/generate", suite.ctrl.Generate).Methods(http.MethodPost)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/generate", bytes.NewBuffer([]byte(body)))
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":false,"error":{"message":"failed to save shortened URL in URL store"}}`
	suite.Equal(http.StatusInternalServerError, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func (suite *GenerateURLControllerTestSuite) TestGenerateReturnServerErrorDueToHashCollisionWhenRetriesExhausted() {
	body := `{"original_url": "https://www.google.com?q=infracloud"}`
	deserializedRequest := contracts.GenerateRequest{OriginalURL: "https://www.google.com?q=infracloud"}

	suite.urlGenerationService.On("Generate", deserializedRequest).Return(contracts.GenerateResponse{}, errs.ErrURLGenerationHashCollisionURL)
	suite.router.HandleFunc("/v1/generate", suite.ctrl.Generate).Methods(http.MethodPost)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/generate", bytes.NewBuffer([]byte(body)))
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":false,"error":{"message":"hash collision has occurred: same hash generated for different url, please retry"}}`
	suite.Equal(http.StatusInternalServerError, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func (suite *GenerateURLControllerTestSuite) TestGenerateReturnSuccessResponseWithRetryDueToHashCollision() {
	body := `{"original_url": "https://www.google.com?q=infracloud"}`
	deserializedRequest := contracts.GenerateRequest{OriginalURL: "https://www.google.com?q=infracloud"}
	response := contracts.GenerateResponse{TinyURL: "http://localhost:8081/01LY7VK"}

	suite.urlGenerationService.On("Generate", deserializedRequest).Return(contracts.GenerateResponse{}, errs.ErrURLGenerationHashCollisionURL).Once()
	suite.urlGenerationService.On("Generate", deserializedRequest).Return(response, nil).Once()
	suite.router.HandleFunc("/v1/generate", suite.ctrl.Generate).Methods(http.MethodPost)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/generate", bytes.NewBuffer([]byte(body)))
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":true,"data":{"tiny_url":"http://localhost:8081/01LY7VK"}}`
	suite.Equal(http.StatusCreated, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func (suite *GenerateURLControllerTestSuite) TestGenerateReturnSuccessResponse() {
	body := `{"original_url": "https://www.google.com?q=infracloud"}`
	deserializedRequest := contracts.GenerateRequest{OriginalURL: "https://www.google.com?q=infracloud"}
	response := contracts.GenerateResponse{TinyURL: "http://localhost:8081/01LY7VK"}

	suite.urlGenerationService.On("Generate", deserializedRequest).Return(response, nil)
	suite.router.HandleFunc("/v1/generate", suite.ctrl.Generate).Methods(http.MethodPost)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/generate", bytes.NewBuffer([]byte(body)))
	suite.router.ServeHTTP(recorder, request)

	responseBytes, _ := io.ReadAll(recorder.Body)
	expectedResponse := `{"success":true,"data":{"tiny_url":"http://localhost:8081/01LY7VK"}}`
	suite.Equal(http.StatusCreated, recorder.Code)
	suite.Equal(expectedResponse, string(responseBytes))
}

func TestGenerateURLControllerTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateURLControllerTestSuite))
}
