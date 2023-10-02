package controllers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/isomnath/tiny-url/contracts"
)

type URLRedirectionController struct {
	urlRedirectionService urlRedirectionService
}

// Redirect - Redirect to original URL
// @Summary Redirect
// @Description Redirect to original URL
// @Tags Functional
// @Produce json
// @Param tiny-url path string true "01LY7VK"
// @Router /{tiny-url} [get]
func (ctrl *URLRedirectionController) Redirect(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	args := mux.Vars(r)
	shortenedRep := args["tiny-url"]

	originalURL, err := ctrl.urlRedirectionService.Redirect(ctx, shortenedRep)
	if err != nil {
		contracts.ErrorResponse(rw, err, contracts.ErrorInternalServer)
		return
	}

	http.Redirect(rw, r, originalURL, http.StatusOK)
}

func NewURLRedirectionController(urlRedirectionService urlRedirectionService) *URLRedirectionController {
	return &URLRedirectionController{
		urlRedirectionService: urlRedirectionService,
	}
}
