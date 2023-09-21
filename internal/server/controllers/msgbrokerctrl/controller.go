package msgbrokerctrl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/msghelper"
)

const packageName = "msgbrokerctrl"

type Controller struct {
	//INPUT channels.
	chIn <-chan msgbroker.Message

	//OUTPUT channels.
	chError chan<- msgbroker.Message
	chOut   chan<- datastructs.ExtraDataFilling

	log logger.BaseLogger
}

func Create(chIn <-chan msgbroker.Message, chError chan<- msgbroker.Message, chPartialPersonData chan<- datastructs.ExtraDataFilling,
	log logger.BaseLogger) *Controller {
	return &Controller{chIn: chIn,
		chError: chError,
		chOut:   chPartialPersonData,
		log:     log,
	}
}

func (c *Controller) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(c.chError)
			close(c.chOut)
			return

		case msgIn, ok := <-c.chIn:
			if !ok {
				close(c.chError)
				close(c.chOut)
				return
			}

			if err := c.handleMessage(msgIn); err != nil {
				if err := msghelper.PutErrorMessageInChan(c.chError, &msgIn, "error-in-incoming-msg", err); err != nil {
					c.log.Info("["+packageName+":Controller:Run] put message in error chan: %w", err)
				}
			}
		}
	}
}

func (c *Controller) handleMessage(msg msgbroker.Message) error {
	personData := datastructs.PersonData{}
	if err := json.Unmarshal(msg.Value, &personData); err != nil {
		return fmt.Errorf("msg JSON unmarshaling: %w", err)
	}

	_, err := msghelper.IsMessageValid(personData)
	if err != nil {
		return fmt.Errorf("input messsage is incorrect: %w", err)
	}

	c.chOut <- datastructs.ExtraDataFilling{
		Raw:  msg,
		Data: personData,
	}
	c.log.Info("["+packageName+":Controller:handleMessage] person data from msg has been prepared to fill extra data: %v", personData)
	return nil
}
