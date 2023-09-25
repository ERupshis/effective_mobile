package msghelper // import "github.com/erupshis/effective_mobile/internal/server/helpers/msghelper"

Package msghelper provides functions for kafka's messages validation, creation
and handling error messages.

FUNCTIONS

func CreateErrorMessage(originalMsg []byte, err error) ([]byte, error)
    CreateErrorMessage creates error message for kafka.

func IsMessageValid(personData datastructs.PersonData) (bool, error)
    IsMessageValid validates incoming message from kafka(name and surname).

func PutErrorMessageInChan(chError chan<- msgbroker.Message, msg *msgbroker.Message, errMsgKey string, err error) error
    PutErrorMessageInChan creates error message for kafka with message key and
    body and puts it in channel for sending into kafka's topic.

