package httpctrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/httpctrl"

Package httpctrl rest api controller. Handles http requests

TYPES

type Controller struct {
	// Has unexported fields.
}
    Controller struct of controller. Receives CRUD requests, validate them and
    forwards to storage if requests are correct.

func Create(strg storage.BaseStorage, log logger.BaseLogger) *Controller
    Create creates controller.

func (c *Controller) Route() *chi.Mux
    Route handlers router.

