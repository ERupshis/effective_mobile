package storage

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

type Storage interface {
	SavePersonData(ctx context.Context, data *datastructs.PersonData) error
}
