package qraphqlctrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/qraphqlctrl"

Package qraphqlctrl provides GraphQL handling.

TYPES

type Controller struct {
	// Has unexported fields.
}

func Create(strg storage.BaseStorage, log logger.BaseLogger) *Controller
    Create returns controller.

func (c *Controller) Route() *chi.Mux
    Route generates routing for controller.

