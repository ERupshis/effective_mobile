package errorsctrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/errorsctrl"


TYPES

type Controller struct {
	// Has unexported fields.
}

func Create(chansIn []<-chan msgbroker.Message, chanOut chan<- msgbroker.Message, log logger.BaseLogger) *Controller

func (c *Controller) Run(ctx context.Context)

