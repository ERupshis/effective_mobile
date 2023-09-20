package storage

import (
	"github.com/erupshis/effective_mobile/internal/datastructs"
)

type ram struct {
	data map[int]datastructs.PersonData
}

func CreateRamStorage() Storage {
	storage := ram{}
	storage.data = make(map[int]datastructs.PersonData)
	return &storage
}

func (r *ram) SavePersonData(data *datastructs.PersonData) error {
	r.data[len(r.data)] = *data
	return nil
}
