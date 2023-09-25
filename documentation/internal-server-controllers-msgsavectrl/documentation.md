package msgsavectrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/msgsavectrl"

Package msgsavectrl receives prepared person data ins storage.

TYPES

type Controller struct {
	// Has unexported fields.
}
    Controller struct of controller. Receives messages from extra data ctrl, and
    adds it in storage. In case or error sends error message in error channel.

func Create(chIn <-chan datastructs.ExtraDataFilling, chError chan<- msgbroker.Message, strg storage.BaseStorage, log logger.BaseLogger) *Controller
    Create creates controller.

func (c *Controller) Run(ctx context.Context)
    Run main goroutine listens input channel and processes it.

