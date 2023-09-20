package msgbrokercontroller

import (
	"context"
	"encoding/json"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/storage"
)

// TODO: should handle messages and write it in storage.
// TODO: in case of error should interrupt and send error message in topic.
type Controller struct {
	chIn    <-chan msgbroker.Message
	chError chan<- msgbroker.Message

	strg storage.Storage
	log  logger.BaseLogger
}

func Create(chIn <-chan msgbroker.Message, chError chan<- msgbroker.Message, strg storage.Storage, log logger.BaseLogger) *Controller {
	return &Controller{chIn: chIn, chError: chError, strg: strg, log: log}
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

			//TODO: validation, adders for data

			personData := datastructs.PersonData{}
			if err := json.Unmarshal(msgIn.Value, &personData); err != nil {
				c.log.Info("failed to parse JSON message: %s, error: %v", msgIn.Value, err)
				//TODO: need to create error msg.
				continue
			}

			if err := c.strg.SavePersonData(&personData); err != nil {
				c.log.Info("storage save value fail. message: %s, error: %v", msgIn.Value, err)
				//TODO: need to create error msg.
			}
			c.log.Info("person data from message saved: %v\n", personData)
		}
	}
}
