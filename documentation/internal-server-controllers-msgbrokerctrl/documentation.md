package msgbrokerctrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/msgbrokerctrl"

Package msgbrokerctrl receives raw messages from broker

TYPES

type Controller struct {
	// Has unexported fields.
}
    Controller struct of controller. Receives messages from broker, validates
    them and parses into person's data struct with preliminary validation.
    put into output channel in case of success. otherwise put error in error's
    channel.

func Create(chIn <-chan msgbroker.Message, chError chan<- msgbroker.Message, chPartialPersonData chan<- datastructs.ExtraDataFilling, log logger.BaseLogger) *Controller
    Create creates controller.

func (c *Controller) Run(ctx context.Context)
    Run main goroutine listens input channel and processes it.

