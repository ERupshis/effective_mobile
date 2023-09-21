package msgbrokercontroller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/storage"
)

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

			if err := c.handleMessage(ctx, msgIn); err != nil {
				msgErr, err := c.createErrorMessage(msgIn.Value, err)
				if err != nil {
					c.log.Info("create error message: %w", err)
					continue
				}

				c.chError <- msgbroker.Message{
					Key:   []byte("error-in-incoming-message"),
					Value: msgErr,
				}
			}
		}
	}
}

func (c *Controller) handleMessage(ctx context.Context, message msgbroker.Message) error {
	personData := datastructs.PersonData{}
	if err := json.Unmarshal(message.Value, &personData); err != nil {
		return fmt.Errorf("message JSON unmarshaling: %w", err)
	}

	_, err := isMessageValid(personData)
	if err != nil {
		return fmt.Errorf("input messsage is incorrect: %w", err)
	}

	if err := c.strg.SavePersonData(ctx, &personData); err != nil {
		return fmt.Errorf("storage save value fail: %w", err)
	}

	c.log.Info("person data from message saved: %v\n", personData)
	return nil
}

func (c *Controller) createErrorMessage(originalMsg []byte, err error) ([]byte, error) {
	msgError, errMarshaling := json.Marshal(
		datastructs.ErrorMessage{
			OriginalMessage: string(originalMsg),
			Error:           err.Error(),
		})

	if errMarshaling != nil {
		return []byte{}, fmt.Errorf("JSON marshaling error: %w", errMarshaling)
	}

	return msgError, nil
}

func isMessageValid(personData datastructs.PersonData) (bool, error) {
	var err error
	if personData.Surname == "" {
		err = fmt.Errorf("surname is/are empty")
	}

	if personData.Name == "" {
		if err == nil {
			err = fmt.Errorf("name is empty")
		} else {
			err = fmt.Errorf("name and %v", err)
		}
	}

	return err == nil, err
}
