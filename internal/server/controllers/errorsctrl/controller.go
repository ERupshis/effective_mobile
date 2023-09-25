package errorsctrl

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

const packageName = "errorsctrl"

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

	for {
		if _, ok := <-ctx.Done(); !ok {
			c.log.Info("[%s:Controller:fanInMessages] close output channel(context is over).", packageName)
			stopCh <- struct{}{}
			close(c.chOut)
			c.chOut = nil
			return
		}
	}
}

func (c *Controller) fanInMessages(stopCh <-chan struct{}) {
	c.log.Info("[%s:Controller:fanInMessages] starting collectors channels handling in goroutines. count: %d", packageName, len(c.chansIn))
	for _, chIn := range c.chansIn {
		go func(stopCh <-chan struct{}, chIn <-chan msgbroker.Message, chOut chan<- msgbroker.Message) {
			for errMsg := range chIn {
				select {
				case <-stopCh:
					return
				default:
					if c.chOut == nil {
						c.log.Info("[%s:Controller:fanInMessages] failed to add message. Output chan is closed: %d", packageName)
						return
					}
					c.chOut <- errMsg
				}
			}
			c.log.Info("[%s:Controller:fanInMessages] close fanIn chan. Total count: %d", packageName, len(c.chansIn))

		}(stopCh, chIn, c.chOut)
	}
}
