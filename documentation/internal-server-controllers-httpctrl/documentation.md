package httpctrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/httpctrl"


TYPES

type Controller struct {
	// Has unexported fields.
}

func Create(strg storage.BaseStorage, log logger.BaseLogger) *Controller

func (c *Controller) Route() *chi.Mux

