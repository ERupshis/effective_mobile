package msgbrokerctrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/msgbrokerctrl"


TYPES

type Controller struct {
	// Has unexported fields.
}

func Create(chIn <-chan msgbroker.Message, chError chan<- msgbroker.Message, chPartialPersonData chan<- datastructs.ExtraDataFilling,
	log logger.BaseLogger) *Controller

func (c *Controller) Run(ctx context.Context)

