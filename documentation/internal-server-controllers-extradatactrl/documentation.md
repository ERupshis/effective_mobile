package extradatactrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/extradatactrl"

Package extradatactr remote api calls controller.

TYPES

type Controller struct {
	// Has unexported fields.
}
    Controller struct of controller. Receives data which should be filled with
    additional data, taken from remote services. In case of success puts person
    data with extra data to further operating. Otherwise - generates error
    message in error channel.

func Create(chIn <-chan datastructs.ExtraDataFilling, chOut chan<- datastructs.ExtraDataFilling, chError chan<- msgbroker.Message,
	client client.BaseClient, log logger.BaseLogger) *Controller
    Create creates controller.

func (c *Controller) Run(ctx context.Context)
    Run main goroutine listens input channel and processes it.

