package msghelper // import "github.com/erupshis/effective_mobile/internal/server/helpers/msghelper"

Package msghelper provides functions for kafka's messages validation, creation
and handling error messages.

func CreateErrorMessage(originalMsg []byte, err error) ([]byte, error)
func IsMessageValid(personData datastructs.PersonData) (bool, error)
func PutErrorMessageInChan(chError chan<- msgbroker.Message, msg *msgbroker.Message, errMsgKey string, ...) error
