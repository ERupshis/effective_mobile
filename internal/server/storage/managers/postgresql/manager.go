// Package postgresql postgresql handling PostgreSQL database.
package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/retryer"
	"github.com/erupshis/effective_mobile/internal/server/config"
	"github.com/erupshis/effective_mobile/internal/server/storage/managers"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// postgresDB storageManager implementation for PostgreSQL. Consist of database and QueriesHandler.
// Request to database are synchronized by sync.RWMutex. All requests is done on united transaction. Multi insert/update/delete is not supported at the moment.
type postgresDB struct {
	database *sql.DB
	handler  QueriesHandler

	log logger.BaseLogger
	mu  sync.RWMutex
}

// CreatePostgreDB creates manager implementation. Supports migrations and check connection to database.
func CreatePostgreDB(ctx context.Context, cfg config.Config, queriesHandler QueriesHandler, log logger.BaseLogger) (managers.BaseStorageManager, error) {
	log.Info("[CreatePostgreDB] open database with settings: '%s'", cfg.DatabaseDSN)
	createDatabaseError := "create db: %w"
	database, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	manager := &postgresDB{
		database: database,
		handler:  queriesHandler,
		log:      log,
	}

	if _, err = manager.CheckConnection(ctx); err != nil {
		return nil, fmt.Errorf(createDatabaseError, err)
	}

	log.Info("[CreatePostgreDB] successful")
	return manager, nil
}

// CheckConnection checks connection to database.
func (p *postgresDB) CheckConnection(ctx context.Context) (bool, error) {
	exec := func(context context.Context) (int64, []byte, error) {
		return 0, []byte{}, p.database.PingContext(context)
	}
	_, _, err := retryer.RetryCallWithTimeout(ctx, p.log, nil, databaseErrorsToRetry, exec)
	if err != nil {
		return false, fmt.Errorf("check connection: %w", err)
	}
	return true, nil
}

// Close closes database.
func (p *postgresDB) Close() error {
	return p.database.Close()
}

// AddPerson add person in database. Creates new values in referenced tables.
func (p *postgresDB) AddPerson(ctx context.Context, data *datastructs.PersonData) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.log.Info("[postgresDB:AddPerson] start transaction")
	errorMessage := "add person in db: %w"
	tx, err := p.database.BeginTx(ctx, nil)
	if err != nil {
		return -1, fmt.Errorf(errorMessage, err)
	}

	genderId, err := p.handler.GetAdditionalId(ctx, tx, data.Gender, GendersTable)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return -1, fmt.Errorf(errorMessage, err)
	}

	countryId, err := p.handler.GetAdditionalId(ctx, tx, data.Country, CountriesTable)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return -1, fmt.Errorf(errorMessage, err)
	}

	newPersonId, err := p.handler.InsertPerson(ctx, tx, data, genderId, countryId)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return -1, fmt.Errorf(errorMessage, err)
	}

	err = tx.Commit()
	if err != nil {
		return -1, fmt.Errorf(errorMessage, err)
	}

	p.log.Info("[postgresDB:AddPerson] transaction successful")
	return newPersonId, nil
}

// SelectPersons returns persons from database satisfying to filters. Supports result pagination.
func (p *postgresDB) SelectPersons(ctx context.Context, filters map[string]interface{}, pageNum int64, pageSize int64) ([]datastructs.PersonData, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	p.log.Info("[postgresDB:SelectPersons] start transaction")
	errorMessage := "select persons in db: %w"
	tx, err := p.database.BeginTx(ctx, nil)
	if err != nil {
		return []datastructs.PersonData{}, fmt.Errorf(errorMessage, err)
	}

	persons, err := p.handler.SelectPersons(ctx, tx, filters, pageNum, pageSize)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return []datastructs.PersonData{}, fmt.Errorf(errorMessage, err)
	}

	err = tx.Commit()
	if err != nil {
		return []datastructs.PersonData{}, fmt.Errorf(errorMessage, err)
	}

	p.log.Info("[postgresDB:SelectPersons] transaction successful")
	return persons, nil
}

// DeletePersonById deletes person by id.
func (p *postgresDB) DeletePersonById(ctx context.Context, personId int64) (int64, error) {
	//TODO: avoid real deletion from DB. Need to add new attr 'isDeleted' and mark on deleted elements.
	p.mu.Lock()
	defer p.mu.Unlock()

	p.log.Info("[postgresDB:DeletePerson] start transaction")
	errorMessage := "delete person in db: %w"
	tx, err := p.database.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf(errorMessage, err)
	}

	affectedCount, err := p.handler.DeletePerson(ctx, tx, personId)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return 0, fmt.Errorf(errorMessage, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf(errorMessage, err)
	}

	p.log.Info("[postgresDB:DeletePerson] transaction successful")

	return affectedCount, nil
}

// UpdatePersonById updates person by id.
func (p *postgresDB) UpdatePersonById(ctx context.Context, id int64, values map[string]interface{}) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.log.Info("[postgresDB:UpdatePartiallyPersonById] start transaction")
	errorMessage := "update person partially in db: %w"
	tx, err := p.database.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf(errorMessage, err)
	}

	err = p.replaceRefValues(ctx, tx, values)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return 0, fmt.Errorf(errorMessage, err)
	}

	affectedCount, err := p.handler.UpdatePartialPersonById(ctx, tx, id, values)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return 0, fmt.Errorf(errorMessage, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf(errorMessage, err)
	}

	p.log.Info("[postgresDB:UpdatePartiallyPersonById] transaction successful")
	return affectedCount, nil
}

// replaceRefValues makes substitution of referenced values in request to foreign key name.
func (p *postgresDB) replaceRefValues(ctx context.Context, tx *sql.Tx, values map[string]interface{}) error {
	valuesToReplace := []struct {
		name  string
		table string
	}{
		{
			name:  "gender",
			table: GendersTable,
		},
		{
			name:  "country",
			table: CountriesTable,
		},
	}

	errorMessage := "replace reference fields in db: %w"
	for _, value := range valuesToReplace {
		if incomingVal, ok := values[value.name]; ok {
			incomingVal := helpers.InterfaceToString(incomingVal)

			valId, err := p.handler.GetAdditionalId(ctx, tx, incomingVal, value.table)
			if err != nil {
				helpers.ExecuteWithLogError(tx.Rollback, p.log)
				return fmt.Errorf(errorMessage, err)
			}

			delete(values, value.name)
			values[value.name+"_id"] = strconv.FormatInt(valId, 10)
		}
	}

	return nil
}
