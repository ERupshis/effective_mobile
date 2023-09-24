// Package msghelper provides functions for kafka's messages validation, creation and handling error messages.
package msghelper

import (
	"encoding/json"
	"fmt"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

// PutErrorMessageInChan creates error message for kafka with message key and body and puts it in channel for sending into kafka's topic.
func PutErrorMessageInChan(chError chan<- msgbroker.Message, msg *msgbroker.Message, errMsgKey string, err error) error {
	if chError == nil {
		return fmt.Errorf("channel is nil")
	}

	if msg == nil {
		return fmt.Errorf("msg is nil")
	}

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

// CreateErrorMessage creates error message for kafka.
func CreateErrorMessage(originalMsg []byte, err error) ([]byte, error) {
	if err == nil {
		return []byte{}, fmt.Errorf("error is nil")
	}

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

// IsMessageValid validates incoming message from kafka(name and surname).
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
