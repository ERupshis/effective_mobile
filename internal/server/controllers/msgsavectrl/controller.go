package msgsavectrl

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/storage"
)

const packageName = "msgsavectrl"

type Controller struct {
	//INPUT channels.
	chIn <-chan datastructs.ExtraDataFilling

	//OUTPUT channels.
	chError chan<- msgbroker.Message

	log  logger.BaseLogger
	strg storage.Storage
}

func Create(chIn <-chan datastructs.ExtraDataFilling, chError chan<- msgbroker.Message, strg storage.Storage, log logger.BaseLogger) *Controller {
	return &Controller{
		chIn:    chIn,
		chError: chError,
		strg:    strg,
		log:     log,
	}
}

func (c *Controller) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(c.chError)
			return

		case msgIn, ok := <-c.chIn:
			if !ok {
				close(c.chError)
				return
			}

			if err := c.strg.AddPersonData(ctx, &msgIn.Data); err != nil {
				c.log.Info("["+packageName+":Controller:Run] failed to add message in storage: %w", err)
				continue
			}

			c.log.Info("["+packageName+":Controller:Run] person data has been saved in storage: %v", msgIn.Data)
		}
	}
}
