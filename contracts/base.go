package contracts

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/isomnath/tiny-url/log"
)

type Err struct {
	Message string `json:"message,omitempty"`
}

type BaseResponse struct {
	Success bool        `json:"success"`
	Error   *Err        `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ReadRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Log.Errorf("failed to read request body with error: %v", err)
		return nil, errors.New("failed to read request body")
	}
	return body, nil
}

func UnmarshalRequest(r *http.Request, destination interface{}) error {
	body, err := ReadRequestBody(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, destination)
	if err != nil {
		log.Log.Errorf("failed to deserialize json request body to destination interface with error: %v", err)
		return errors.New("failed to deserialize json request body to destination interface")
	}
	return nil
}

func SuccessResponse(rw http.ResponseWriter, responseData interface{}, status string) {
	successDetails := successObjects[status]
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(successDetails.status)
	response := BaseResponse{
		Success: true,
		Data:    responseData,
	}

	responseJSON, _ := json.Marshal(response)
	rw.Write(responseJSON)

	return
}

func ErrorResponse(rw http.ResponseWriter, er error, status string) {
	httpStatus := errorObjects[status].status
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(httpStatus)
	response := BaseResponse{
		Success: false,
		Error:   &Err{Message: er.Error()},
	}

	responseJSON, _ := json.Marshal(response)
	rw.Write(responseJSON)
	return
}
