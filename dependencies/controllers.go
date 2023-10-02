package dependencies

import (
	"github.com/isomnath/tiny-url/controllers"
)

type Controllers struct {
	AnalyticsController      *controllers.AnalyticsController
	GenerateURLController    *controllers.GenerateURLController
	URLRedirectionController *controllers.URLRedirectionController
}

func initializeControllers(services *Services) *Controllers {
	return &Controllers{
		AnalyticsController:      controllers.NewAnalyticsController(services.AnalyticsService),
		GenerateURLController:    controllers.NewGenerateURLController(services.URLGenerationService),
		URLRedirectionController: controllers.NewURLRedirectionController(services.URLRedirectionService),
	}
}
