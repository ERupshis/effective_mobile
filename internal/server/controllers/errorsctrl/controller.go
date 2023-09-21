package errorsctrl

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

type Controller struct {
	//INPUT channels.
	chansIn []<-chan msgbroker.Message

	//OUTPUT channels.
	chOut chan<- msgbroker.Message

	log logger.BaseLogger
}

func Create(chansIn []<-chan msgbroker.Message, chanOut chan<- msgbroker.Message, log logger.BaseLogger) *Controller {
	return &Controller{
		chansIn: chansIn,
		chOut:   chanOut,
		log:     log,
	}
}

func (c *Controller) Run(ctx context.Context) {
	stopCh := make(chan struct{})
	c.fanInMessages(stopCh)

	for range ctx.Done() {
		stopCh <- struct{}{}
		close(c.chOut)
		return
	}

}

func (c *Controller) fanInMessages(stopCh <-chan struct{}) {
	for _, chIn := range c.chansIn {
		go func(stopCh <-chan struct{}, chIn <-chan msgbroker.Message, chOut chan<- msgbroker.Message) {
			for errMsg := range chIn {
				select {
				case <-stopCh:
					return
				default:
					c.chOut <- errMsg
				}
			}
		}(stopCh, chIn, c.chOut)
	}
}
