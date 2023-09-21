package msghelper

import (
	"encoding/json"
	"fmt"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

type Controller struct {
	chError chan<- msgbroker.Message

	log logger.BaseLogger
}

func PutErrorMessageInChan(chError chan<- msgbroker.Message, msg *msgbroker.Message, errMsgKey string, err error) error {
	msgErr, err := CreateErrorMessage(msg.Value, err)
	if err != nil {
		return fmt.Errorf("create error message: %w", err)
	}

	chError <- msgbroker.Message{
		Key:   []byte(errMsgKey),
		Value: msgErr,
	}

	return nil
}

func CreateErrorMessage(originalMsg []byte, err error) ([]byte, error) {
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

func IsMessageValid(personData datastructs.PersonData) (bool, error) {
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
