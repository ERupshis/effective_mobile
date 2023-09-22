package postgrequeries

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/retryer"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
)

const (
	SchemaName     = "persons_data"
	PersonsTable   = "persons"
	GendersTable   = "genders"
	CountriesTable = "countries"
)

var FieldsInPersonsTable = []string{"name", "surname", "patronymic", "age", "gender_id", "country_id"}

var DatabaseErrorsToRetry = []error{
	errors.New(pgerrcode.UniqueViolation),
	errors.New(pgerrcode.ConnectionException),
	errors.New(pgerrcode.ConnectionDoesNotExist),
	errors.New(pgerrcode.ConnectionFailure),
	errors.New(pgerrcode.SQLClientUnableToEstablishSQLConnection),
	errors.New(pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection),
	errors.New(pgerrcode.TransactionResolutionUnknown),
	errors.New(pgerrcode.ProtocolViolation),
}

func GetTableFullName(table string) string {
	return SchemaName + "." + table
}

func CreateValuesFormForStmt(count int) string {
	if count <= 0 {
		return ""
	} else if count == 1 {
		return "(?)"
	} else {
		return fmt.Sprintf("(?%s)", strings.Repeat(", ?", count-1))
	}
}

type QueriesHandler struct {
	log logger.BaseLogger
}

func CreateHandler(log logger.BaseLogger) QueriesHandler {
	return QueriesHandler{log: log}
}

func (q *QueriesHandler) InsertPerson(ctx context.Context, tx *sql.Tx, personData *datastructs.PersonData, genderId int64, countryId int64) error {
	errorMsg := fmt.Sprintf("insert person '%v' in '%s': %w", personData, PersonsTable)

	stmt, err := createInsertPersonStmt(ctx, tx)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	defer helpers.ExecuteWithLogError(stmt.Close, q.log)

	query := func(context context.Context) error {
		_, err = stmt.ExecContext(
			context,
			personData.Name,
			personData.Surname,
			personData.Patronymic,
			personData.Age,
			genderId,
			countryId,
		)
		return err
	}
	err = retryer.RetryCallWithTimeoutErrorOnly(ctx, q.log, []int{1, 1, 3}, DatabaseErrorsToRetry, query)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}

func createInsertPersonStmt(ctx context.Context, tx *sql.Tx) (*sql.Stmt, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	psqlInsert, _, err := psql.Insert(GetTableFullName(PersonsTable)).
		Columns(FieldsInPersonsTable...).
		Values(make([]interface{}, len(FieldsInPersonsTable))...).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("squirrel sql insert statement for '"+GetTableFullName(PersonsTable)+"': %w", err)

	}
	return tx.PrepareContext(ctx, psqlInsert)
}

func (q *QueriesHandler) GetAdditionalId(ctx context.Context, tx *sql.Tx, name string, table string) (int64, error) {
	errorMsg := fmt.Sprintf("get additional id for '%s' in '%s': %w", name, table)

	stmt, err := createSelectAdditionalIdStmt(ctx, tx, name, table)
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}
	defer helpers.ExecuteWithLogError(stmt.Close, q.log)

	var id int64
	query := func(context context.Context) error {
		return stmt.QueryRowContext(ctx, name).Scan(&id)
	}
	err = retryer.RetryCallWithTimeoutErrorOnly(ctx, q.log, []int{1, 1, 3}, DatabaseErrorsToRetry, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			id, err = q.InsertAdditionalId(ctx, tx, name, table)
			if err != nil {
				return 0, fmt.Errorf(errorMsg, err)
			}
		} else {
			return 0, fmt.Errorf(errorMsg, err)
		}
	}

	return id, nil
}

func createSelectAdditionalIdStmt(ctx context.Context, tx *sql.Tx, name string, table string) (*sql.Stmt, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	psqlSelect, _, err := psql.Select("id").
		From(GetTableFullName(table)).
		Where(sq.Eq{"name": name}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("squirrel sql select statement for '"+GetTableFullName(table)+"': %w", err)

	}
	return tx.PrepareContext(ctx, psqlSelect)
}

func (q *QueriesHandler) InsertAdditionalId(ctx context.Context, tx *sql.Tx, name string, table string) (int64, error) {
	errorMsg := fmt.Sprintf("insert additional value for '%s' in '%s': %w", name, table)

	stmt, err := createInsertAdditionalIdStmt(ctx, tx, name, table)
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}
	defer helpers.ExecuteWithLogError(stmt.Close, q.log)

	var id int64
	query := func(context context.Context) error {
		return stmt.QueryRowContext(ctx, name).Scan(&id)
	}
	err = retryer.RetryCallWithTimeoutErrorOnly(ctx, q.log, []int{1, 1, 3}, DatabaseErrorsToRetry, query)
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}

	return id, nil
}

func createInsertAdditionalIdStmt(ctx context.Context, tx *sql.Tx, name string, table string) (*sql.Stmt, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	psqlInsert, _, err := psql.Insert(GetTableFullName(table)).
		Columns("name").
		Values(name).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("squirrel sql insert statement for '"+GetTableFullName(table)+"': %w", err)

	}
	return tx.PrepareContext(ctx, psqlInsert+"RETURNING id")
}
