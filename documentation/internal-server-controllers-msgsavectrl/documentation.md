package msgsavectrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/msgsavectrl"


TYPES

type Controller struct {
	// Has unexported fields.
}

func Create(chIn <-chan datastructs.ExtraDataFilling, chError chan<- msgbroker.Message, strg storage.BaseStorage, log logger.BaseLogger) *Controller

func (c *Controller) Run(ctx context.Context)

