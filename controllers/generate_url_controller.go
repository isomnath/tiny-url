package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/isomnath/tiny-url/contracts"
	"github.com/isomnath/tiny-url/errs"
)

type GenerateURLController struct {
	urlGenerationService urlGenerationService
}

// Generate - Generate Tiny URL
// @Summary Generate
// @Description Generate Tiny URL
// @Tags Functional
// @Accept json
// @Produce json
// @Router /v1/generate [post]
func (ctrl *GenerateURLController) Generate(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request contracts.GenerateRequest
	err := contracts.UnmarshalRequest(r, &request)
	if err != nil {
		contracts.ErrorResponse(rw, err, contracts.ErrorBadRequest)
		return
	}

	if request.OriginalURL == "" {
		contracts.ErrorResponse(rw, fmt.Errorf("empty url"), contracts.ErrorBadRequest)
		return
	}

	_, err = url.Parse(request.OriginalURL)
	if err != nil {
		contracts.ErrorResponse(rw, fmt.Errorf("invalid url"), contracts.ErrorBadRequest)
		return
	}

	var response contracts.GenerateResponse
	var serverErr error
	for i := 0; i < 3; i++ {
		response, serverErr = ctrl.urlGenerationService.Generate(ctx, request)
		if serverErr != nil {
			if serverErr == errs.ErrURLGenerationHashCollisionURL {
				continue
			} else {
				break
			}
		}
		if err == nil {
			break
		}
	}
	if serverErr != nil {
		contracts.ErrorResponse(rw, serverErr, contracts.ErrorInternalServer)
		return
	}
	contracts.SuccessResponse(rw, response, contracts.SuccessCreated)
	return
}

func NewGenerateURLController(urlGenerationService urlGenerationService) *GenerateURLController {
	return &GenerateURLController{
		urlGenerationService: urlGenerationService,
	}
}
