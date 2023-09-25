package extradatactrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/extradatactrl"


TYPES

type Controller struct {
	// Has unexported fields.
}

func Create(chIn <-chan datastructs.ExtraDataFilling, chOut chan<- datastructs.ExtraDataFilling, chError chan<- msgbroker.Message,
	client client.BaseClient, log logger.BaseLogger) *Controller

func (c *Controller) Run(ctx context.Context)

