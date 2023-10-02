package controllers

import (
	"net/http"

	"github.com/isomnath/tiny-url/contracts"
)

// Ping - Application Health Check API
// @Summary PING
// @Description Check Application Liveliness
// @Tags Health Check
// @Produce json
// @Success 200 {object} contracts.PingResponse
// @Router /ping [get]
func Ping(rw http.ResponseWriter, r *http.Request) {
	res := contracts.PingResponse{Message: "pong"}
	contracts.SuccessResponse(rw, res, contracts.SuccessOK)
}
