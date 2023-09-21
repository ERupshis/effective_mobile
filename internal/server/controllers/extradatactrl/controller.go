package extradatactrl

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/storage"
)

const packageName = "extradatactrl"

type Controller struct {
	chIn  <-chan datastructs.ExtraDataFilling
	chOut chan<- datastructs.ExtraDataFilling

	log  logger.BaseLogger
	strg storage.Storage
}

func Create(chIn <-chan datastructs.ExtraDataFilling, chOut chan<- datastructs.ExtraDataFilling, strg storage.Storage, log logger.BaseLogger) *Controller {
	return &Controller{
		chIn:  chIn,
		chOut: chOut,
		strg:  strg,
		log:   log,
	}
}

func (c *Controller) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(c.chOut)
			return

		case personDataIn, ok := <-c.chIn:
			if !ok {
				close(c.chOut)
				return
			}

			if err := c.strg.AddPersonData(ctx, &personDataIn.Data); err != nil {
				c.log.Info("storage save value fail: %v", err)
			}

			//TODO: add useful work.
			//TODO: in case of error need to log in error channel.
			c.log.Info("storage save value : %s", personDataIn.Data)
		}
	}
}
