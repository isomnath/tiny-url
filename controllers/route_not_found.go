package controllers

import (
	"fmt"
	"net/http"

	"github.com/isomnath/tiny-url/contracts"
)

// RouteNotFoundHandler - All requests are redirected to this router when a request is received for an unsupported path
func RouteNotFoundHandler(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	contracts.ErrorResponse(rw, fmt.Errorf("route %s not found", path), contracts.ErrorNotFound)
}
