package storage

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
	"github.com/erupshis/effective_mobile/internal/server/storage/postgrequeries"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type postgresDB struct {
	database *sql.DB
	handler  postgrequeries.QueriesHandler

	log logger.BaseLogger
	mu  sync.RWMutex
}

func CreatePostgreDB(ctx context.Context, cfg config.Config, queriesHandler postgrequeries.QueriesHandler, log logger.BaseLogger) (BaseStorageManager, error) {
	log.Info("[storage:CreatePostgreDB] open database with settings: '%s'", cfg.DatabaseDSN)
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

	log.Info("[storage:CreatePostgreDB] successful")
	return manager, nil
}

func (p *postgresDB) CheckConnection(ctx context.Context) (bool, error) {
	exec := func(context context.Context) (int64, []byte, error) {
		return 0, []byte{}, p.database.PingContext(context)
	}
	_, _, err := retryer.RetryCallWithTimeout(ctx, p.log, nil, postgrequeries.DatabaseErrorsToRetry, exec)
	if err != nil {
		return false, fmt.Errorf("check connection: %w", err)
	}
	return true, nil
}

func (p *postgresDB) Close() error {
	return p.database.Close()
}

func (p *postgresDB) AddPerson(ctx context.Context, data *datastructs.PersonData) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.log.Info("[postgresDB:AddPerson] start transaction")
	addPersonError := "add person in db: %w"
	tx, err := p.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(addPersonError, err)
	}

	genderId, err := p.handler.GetAdditionalId(ctx, tx, data.Gender, postgrequeries.GendersTable)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return fmt.Errorf(addPersonError, err)
	}

	countryId, err := p.handler.GetAdditionalId(ctx, tx, data.Country, postgrequeries.CountriesTable)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return fmt.Errorf(addPersonError, err)
	}

	err = p.handler.InsertPerson(ctx, tx, data, genderId, countryId)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return fmt.Errorf(addPersonError, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf(addPersonError, err)
	}

	p.log.Info("[postgresDB:AddPerson] transaction successful")
	return nil
}

func (p *postgresDB) GetPersons(ctx context.Context, filters map[string]string, page int64, pageSize int64) ([]datastructs.PersonData, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	p.log.Info("GetPersons is not implemented")
	return []datastructs.PersonData{}, nil
}

func (p *postgresDB) DeletePersonById(ctx context.Context, personId int64) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.log.Info("[postgresDB:DeletePerson] start transaction")
	deletePersonError := "delete person in db: %w"
	tx, err := p.database.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf(deletePersonError, err)
	}

	affectedCount, err := p.handler.DeletePerson(ctx, tx, personId)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return 0, fmt.Errorf(deletePersonError, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf(deletePersonError, err)
	}

	p.log.Info("[postgresDB:DeletePerson] transaction successful")

	return affectedCount, nil
}

func (p *postgresDB) UpdatePersonById(ctx context.Context, id int64, data *datastructs.PersonData) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.log.Info("[postgresDB:UpdatePersonById] start transaction")
	addPersonError := "add person in db: %w"
	tx, err := p.database.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf(addPersonError, err)
	}

	genderId, err := p.handler.GetAdditionalId(ctx, tx, data.Gender, postgrequeries.GendersTable)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return 0, fmt.Errorf(addPersonError, err)
	}

	countryId, err := p.handler.GetAdditionalId(ctx, tx, data.Country, postgrequeries.CountriesTable)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return 0, fmt.Errorf(addPersonError, err)
	}

	affectedCount, err := p.handler.UpdatePersonById(ctx, tx, id, data, genderId, countryId)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return 0, fmt.Errorf(addPersonError, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf(addPersonError, err)
	}

	p.log.Info("[postgresDB:UpdatePersonById] transaction successful")
	return affectedCount, nil
}

func (p *postgresDB) UpdatePersonByIdPartially(ctx context.Context, id int64, values map[string]interface{}) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.log.Info("[postgresDB:UpdatePartiallyPersonById] start transaction")
	addPersonError := "add person in db: %w"
	tx, err := p.database.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf(addPersonError, err)
	}

	if gender, ok := values["gender"]; ok {
		gender, err := helpers.InterfaceToString(gender)
		if err != nil {
			helpers.ExecuteWithLogError(tx.Rollback, p.log)
			return 0, fmt.Errorf(addPersonError, err)
		}

		genderId, err := p.handler.GetAdditionalId(ctx, tx, gender, postgrequeries.GendersTable)
		if err != nil {
			helpers.ExecuteWithLogError(tx.Rollback, p.log)
			return 0, fmt.Errorf(addPersonError, err)
		}

		delete(values, "gender")
		values["gender_id"] = strconv.FormatInt(genderId, 10)
	}

	if country, ok := values["country"]; ok {
		country, err := helpers.InterfaceToString(country)
		if err != nil {
			helpers.ExecuteWithLogError(tx.Rollback, p.log)
			return 0, fmt.Errorf(addPersonError, err)
		}

		countryId, err := p.handler.GetAdditionalId(ctx, tx, country, postgrequeries.CountriesTable)
		if err != nil {
			helpers.ExecuteWithLogError(tx.Rollback, p.log)
			return 0, fmt.Errorf(addPersonError, err)
		}

		delete(values, "country")
		values["country_id"] = strconv.FormatInt(countryId, 10)
	}

	affectedCount, err := p.handler.UpdatePartialPersonById(ctx, tx, id, values)
	if err != nil {
		helpers.ExecuteWithLogError(tx.Rollback, p.log)
		return 0, fmt.Errorf(addPersonError, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf(addPersonError, err)
	}

	p.log.Info("[postgresDB:UpdatePartiallyPersonById] transaction successful")
	return affectedCount, nil
}

func (p *postgresDB) replaceRefValues(ctx context.Context, tx *sql.Tx, values map[string]interface{}) error {
	valuesToReplace := []struct {
		name  string
		table string
	}{
		{
			name:  "gender",
			table: postgrequeries.GendersTable,
		},
		{
			name:  "country",
			table: postgrequeries.CountriesTable,
		},
	}

	errorMessage := "replace reference fields in db: %w"
	for _, value := range valuesToReplace {
		if incomingVal, ok := values[value.name]; ok {
			incomingVal, err := helpers.InterfaceToString(incomingVal)
			if err != nil {
				helpers.ExecuteWithLogError(tx.Rollback, p.log)
				return fmt.Errorf(errorMessage, err)
			}

			valId, err := p.handler.GetAdditionalId(ctx, tx, incomingVal, value.table)
			if err != nil {
				helpers.ExecuteWithLogError(tx.Rollback, p.log)
				return fmt.Errorf(errorMessage, err)
			}

			delete(values, "gender")
			values["gender_id"] = strconv.FormatInt(valId, 10)
		}
	}

	return nil
}
