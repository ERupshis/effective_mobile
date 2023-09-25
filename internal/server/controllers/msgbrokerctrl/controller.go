// Package msgbrokerctrl receives raw messages from broker
package msgbrokerctrl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/helpers/msghelper"
)

const packageName = "msgbrokerctrl"

// Controller struct of controller. Receives messages from broker, validates them and parses into person's data struct with preliminary validation.
// put into output channel in case of success. otherwise put error in error's channel.
type Controller struct {
	//INPUT channels.
	chIn <-chan msgbroker.Message

	//OUTPUT channels.
	chError chan<- msgbroker.Message
	chOut   chan<- datastructs.ExtraDataFilling

	log logger.BaseLogger
}

// Create creates controller.
func Create(chIn <-chan msgbroker.Message, chError chan<- msgbroker.Message, chPartialPersonData chan<- datastructs.ExtraDataFilling, log logger.BaseLogger) *Controller {
	return &Controller{
		chIn:    chIn,
		chError: chError,
		chOut:   chPartialPersonData,
		log:     log,
	}
}

// Run main goroutine listens input channel and processes it.
func (c *Controller) Run(ctx context.Context) {
	c.log.Info("[%s:Controller:Run] start work", packageName)

	for {
		select {
		case <-ctx.Done():
			c.log.Info("[%s:Controller:Run] stop work due to context was canceled", packageName)
			close(c.chError)
			close(c.chOut)
			return

		case msgIn, ok := <-c.chIn:
			if !ok {
				c.log.Info("[%s:Controller:Run] stop work due to input channel was closed", packageName)
				close(c.chError)
				close(c.chOut)
				return
			}

			personData, err := c.handleMessage(msgIn)
			if err != nil {
				if err := msghelper.PutErrorMessageInChan(c.chError, &msgIn, "error-in-incoming-msg", err); err != nil {
					c.log.Info("[%s:Controller:Run] put message in error chan failed: %w", packageName, err)
				}
				c.log.Info("[%s:Controller:Run] put message in error chan done: %w", packageName, err)
				continue
			}

			c.log.Info("[%s:Controller:handleMessage] person data from msg has been prepared to fill extra data: %v", packageName, personData)
		}
	}
}

// handleMessage message validatior. Unmarshaling from json and validates it for critical fields presense.
func (c *Controller) handleMessage(msg msgbroker.Message) (*datastructs.PersonData, error) {
	personData := datastructs.PersonData{}
	if err := json.Unmarshal(msg.Value, &personData); err != nil {
		return nil, fmt.Errorf("msg JSON unmarshaling: %w", err)
	}

	_, err := msghelper.IsMessageValid(personData)
	if err != nil {
		return nil, fmt.Errorf("input messsage is incorrect: %w", err)
	}

	c.chOut <- datastructs.ExtraDataFilling{
		Raw:  msg,
		Data: personData,
	}
	return &personData, nil
}
