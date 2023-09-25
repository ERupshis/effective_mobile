// Package msgsavectrl receives prepared person data ins storage.
package msgsavectrl

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/helpers/msghelper"
	"github.com/erupshis/effective_mobile/internal/server/storage"
)

const packageName = "msgsavectrl"

// Controller struct of controller. Receives messages from extra data ctrl, and adds it in storage.
// In case or error sends error message in error channel.
type Controller struct {
	//INPUT channels.
	chIn <-chan datastructs.ExtraDataFilling

	//OUTPUT channels.
	chError chan<- msgbroker.Message

	strg storage.BaseStorage
	log  logger.BaseLogger
}

// Create creates controller.
func Create(chIn <-chan datastructs.ExtraDataFilling, chError chan<- msgbroker.Message, strg storage.BaseStorage, log logger.BaseLogger) *Controller {
	return &Controller{
		chIn:    chIn,
		chError: chError,
		strg:    strg,
		log:     log,
	}
}

// Run main goroutine listens input channel and processes it.
func (c *Controller) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.log.Info("[%s:Controller:Run] stop work due to context was canceled", packageName)
			close(c.chError)
			return

		case msgIn, ok := <-c.chIn:
			if !ok {
				c.log.Info("[%s:Controller:Run] stop work due to input channel was closed", packageName)
				close(c.chError)
				return
			}

			newPersonId, err := c.strg.AddPerson(ctx, &msgIn.Data)
			if err != nil {
				if err := msghelper.PutErrorMessageInChan(c.chError, &msgIn.Raw, "error-on-save", err); err != nil {
					c.log.Info("[%s:Controller:Run] put message in error chan failed: %w", packageName, err)
					continue
				}
				c.log.Info("[%s:Controller:Run] failed to add message in storage: %v", packageName, err)
				continue
			}

			msgIn.Data.Id = newPersonId
			c.log.Info("[%s:Controller:Run] person data has been saved in storage: %v", packageName, msgIn.Data)
		}
	}
}
