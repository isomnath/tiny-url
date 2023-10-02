package server

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/isomnath/tiny-url/dependencies"
)

func InitializeRouter(dep *dependencies.Dependencies) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1/generate", dep.Controllers.GenerateURLController.Generate).Methods(http.MethodPost)
	router.HandleFunc("/{tiny-url}", dep.Controllers.URLRedirectionController.Redirect).Methods(http.MethodGet)
	router.HandleFunc("/v1/analytics/domains/highest_transformation", dep.Controllers.AnalyticsController.TopHighTransformedDomains).Methods(http.MethodGet)
	router.HandleFunc("/v1/analytics/domains/highest_traffic", dep.Controllers.AnalyticsController.TopHighTrafficDomains).Methods(http.MethodGet)
	return router
}
