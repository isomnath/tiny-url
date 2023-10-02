package contracts

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isomnath/tiny-url/config"
	"github.com/isomnath/tiny-url/log"

	"github.com/stretchr/testify/suite"
)

type BaseContractTestSuite struct {
	suite.Suite
}

type errReader int

func (suite *BaseContractTestSuite) SetupTest() {
	config.LoadBaseConfig()
	log.Setup()
}

func (suite *BaseContractTestSuite) TestSuccessResponse() {
	type TestData struct {
		Message string `json:"message"`
	}
	rw := httptest.NewRecorder()

	SuccessResponse(rw, TestData{Message: "successful"}, SuccessOK)
	suite.Equal("{\"success\":true,\"data\":{\"message\":\"successful\"}}", rw.Body.String())
	suite.Equal(http.StatusOK, rw.Code)
}

func (suite *BaseContractTestSuite) TestErrorResponse() {
	err := fmt.Errorf("error")
	rw := httptest.NewRecorder()

	ErrorResponse(rw, err, ErrorBadRequest)
	suite.Equal("{\"success\":false,\"error\":{\"message\":\"error\"}}", rw.Body.String())
	suite.Equal(http.StatusBadRequest, rw.Code)
}

func (suite *BaseContractTestSuite) TestUnmarshalRequestSuccessfully() {
	type TestData struct {
		ID   int64  `json:"id"`
		Data string `json:"data"`
	}
	expectedTestData := TestData{
		ID:   123,
		Data: "test",
	}

	var dest TestData
	r, _ := http.NewRequest(http.MethodPost, "/v1/test/route/1", bytes.NewBuffer([]byte(`{"id": 123, "data": "test"}`)))

	err := UnmarshalRequest(r, &dest)
	suite.NoError(err)
	suite.Equal(expectedTestData, dest)
}

func (suite *BaseContractTestSuite) TestUnmarshalRequestJSONUnmarshalError() {
	type TestData struct {
		ID   int64  `json:"id"`
		Data string `json:"data"`
	}

	var dest TestData
	r, _ := http.NewRequest(http.MethodPost, "/v1/test/route/1", bytes.NewBuffer([]byte(`{"id": "123", "data": "test"}`)))

	err := UnmarshalRequest(r, &dest)
	suite.Equal("failed to deserialize json request body to destination interface", err.Error())
}

func (suite *BaseContractTestSuite) TestUnmarshalRequestIOReaderError() {
	type TestData struct {
		ID   int64  `json:"id"`
		Data string `json:"data"`
	}

	var dest TestData
	r, _ := http.NewRequest(http.MethodPost, "/v1/test/route/1", errReader(0))

	err := UnmarshalRequest(r, &dest)
	suite.Equal("failed to read request body", err.Error())
}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestBaseContractTestSuite(t *testing.T) {
	suite.Run(t, new(BaseContractTestSuite))
}
