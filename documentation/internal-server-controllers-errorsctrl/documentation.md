package errorsctrl // import "github.com/erupshis/effective_mobile/internal/server/controllers/errorsctrl"

Package errorsctrl collects errors from other controllers(of kafka messages
handling) in global error channel.

TYPES

type Controller struct {
	// Has unexported fields.
}
    Controller struct of controller. Controls life cycle of global error
    channel.

func Create(chansIn []<-chan msgbroker.Message, chanOut chan<- msgbroker.Message, log logger.BaseLogger) *Controller
    Create creation method.

func (c *Controller) Run(ctx context.Context)
    Run main goroutine method creates collectors for input channels. In case of
    ctx cancel - closes output channel

