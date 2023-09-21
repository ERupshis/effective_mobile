package storage

import (
	"context"
	"sync"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

type ram struct {
	data map[int]datastructs.PersonData
	mu   sync.RWMutex
}

func CreateRamStorage() Storage {
	storage := ram{}
	storage.data = make(map[int]datastructs.PersonData)
	return &storage
}

func (r *ram) SavePersonData(_ context.Context, data *datastructs.PersonData) error {
	r.mu.Lock()
	r.data[len(r.data)] = *data
	r.mu.Unlock()
	return nil
}
