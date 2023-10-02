package controllers

import (
	"net/http"

	"github.com/isomnath/tiny-url/contracts"
)

type AnalyticsController struct {
	analyticsService analyticsService
}

// TopHighTransformedDomains - TopHighTransformedDomains
// @Summary TopHighTransformedDomains
// @Description Fetch Domains With Highest Transformations
// @Tags Functional
// @Produce json
// @Router /v1/analytics/domains/highest_transformation [get]
func (ctrl *AnalyticsController) TopHighTransformedDomains(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	response, err := ctrl.analyticsService.FetchTopDomainsByTransformations(ctx)
	if err != nil {
		contracts.ErrorResponse(rw, err, contracts.ErrorInternalServer)
		return
	}

	contracts.SuccessResponse(rw, response, contracts.SuccessOK)
	return
}

// TopHighTrafficDomains - TopHighTrafficDomains
// @Summary TopHighTrafficDomains
// @Description Fetch Domains With Highest Traffic
// @Tags Functional
// @Produce json
// @Router /v1/analytics/domains/highest_transformation [get]
func (ctrl *AnalyticsController) TopHighTrafficDomains(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	response, err := ctrl.analyticsService.FetchTopDomainsByTraffic(ctx)
	if err != nil {
		contracts.ErrorResponse(rw, err, contracts.ErrorInternalServer)
		return
	}

	contracts.SuccessResponse(rw, response, contracts.SuccessOK)
	return
}

func NewAnalyticsController(analyticsService analyticsService) *AnalyticsController {
	return &AnalyticsController{
		analyticsService: analyticsService,
	}
}
