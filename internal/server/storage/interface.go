package storage

import (
	"github.com/erupshis/effective_mobile/internal/datastructs"
)

type Storage interface {
	SavePersonData(data *datastructs.PersonData) error
}
