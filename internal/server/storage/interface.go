package storage

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

type Storage interface {
	AddPersonData(ctx context.Context, data *datastructs.PersonData) error
}
