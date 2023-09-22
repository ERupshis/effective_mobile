package postgrequeries

import (
	"context"
	"database/sql"
	"fmt"

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

var ColumnsInPersonsTable = []string{"name", "surname", "patronymic", "age", "gender_id", "country_id"}

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

type QueriesHandler struct {
	log logger.BaseLogger
}

func CreateHandler(log logger.BaseLogger) QueriesHandler {
	return QueriesHandler{log: log}
}

func (q *QueriesHandler) InsertPerson(ctx context.Context, tx *sql.Tx, personData *datastructs.PersonData, genderId int64, countryId int64) error {
	errorMsg := fmt.Sprintf("insert person '%v' in '%s'", personData, PersonsTable) + ": %w"

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
		Columns(ColumnsInPersonsTable...).
		Values(make([]interface{}, len(ColumnsInPersonsTable))...).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("squirrel sql insert statement for '"+GetTableFullName(PersonsTable)+"': %w", err)
	}
	return tx.PrepareContext(ctx, psqlInsert)
}

func (q *QueriesHandler) SelectPersons(ctx context.Context, tx *sql.Tx, filters map[string]interface{}, pageNum int64, pageSize int64) ([]datastructs.PersonData, error) {
	errorMsg := fmt.Sprintf("select persons with filter '%v' in '%s'", filters, PersonsTable) + ": %w"

	stmt, err := createSelectPersonsStmt(ctx, tx, filters, pageNum, pageSize)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}
	defer helpers.ExecuteWithLogError(stmt.Close, q.log)

	var valuesToUpdate []interface{}
	for _, val := range filters {
		valuesToUpdate = append(valuesToUpdate, val)
	}

	var rows *sql.Rows
	query := func(context context.Context) error {
		rows, err = stmt.QueryContext(
			context,
			valuesToUpdate...,
		)
		return err
	}
	err = retryer.RetryCallWithTimeoutErrorOnly(ctx, q.log, []int{1, 1, 3}, DatabaseErrorsToRetry, query)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	defer helpers.ExecuteWithLogError(rows.Close, q.log)
	var res []datastructs.PersonData
	for rows.Next() {
		data := datastructs.PersonData{}
		err := rows.Scan(
			&data.Id,
			&data.Name,
			&data.Surname,
			&data.Patronymic,
			&data.Age,
			&data.Gender,
			&data.Country,
		)
		if err != nil {
			return nil, fmt.Errorf("parse db result: %w", err)
		}

		res = append(res, data)
	}

	return res, nil
}

func createSelectPersonsStmt(ctx context.Context, tx *sql.Tx, filters map[string]interface{}, pageNum int64, pageSize int64) (*sql.Stmt, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	gendersJoin := fmt.Sprintf("LEFT JOIN %s ON %[1]s.id = %s.gender_id", GetTableFullName(GendersTable), GetTableFullName(PersonsTable))
	countriesJoin := fmt.Sprintf("LEFT JOIN %s ON %[1]s.id = %s.country_id", GetTableFullName(CountriesTable), GetTableFullName(PersonsTable))
	builder := psql.Select(
		GetTableFullName(PersonsTable)+".id",
		GetTableFullName(PersonsTable)+".name",
		"surname",
		"patronymic",
		"age",
		GetTableFullName(GendersTable)+".name",
		GetTableFullName(CountriesTable)+".name",
	).
		From(GetTableFullName(PersonsTable)).
		JoinClause(gendersJoin).
		JoinClause(countriesJoin)
	if len(filters) != 0 {
		for key := range filters {
			switch key {
			case "name":
				key = GetTableFullName(PersonsTable) + ".name"
			case "gender":
				key = GetTableFullName(GendersTable) + ".name"
			case "country":
				key = GetTableFullName(CountriesTable) + ".name"
			}
			builder = builder.Where(sq.Eq{key: "?"})
		}
	}
	if pageSize != 0 {
		builder = builder.Limit(uint64(pageSize))
		if pageNum != 0 {
			builder = builder.Offset(uint64((pageNum - 1) * pageSize))
		}
	}
	psqlSelect, _, err := builder.ToSql()

	if err != nil {
		return nil, fmt.Errorf("squirrel sql select statement for '"+GetTableFullName(PersonsTable)+"': %w", err)
	}
	return tx.PrepareContext(ctx, psqlSelect)
}

func (q *QueriesHandler) DeletePerson(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	errorMsg := fmt.Sprintf("delete person by id '%v' in '%s", id, PersonsTable) + ": %w"

	stmt, err := createDeletePersonStmt(ctx, tx)
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}
	defer helpers.ExecuteWithLogError(stmt.Close, q.log)

	var result sql.Result
	query := func(context context.Context) error {
		result, err = stmt.ExecContext(context, id)
		return err
	}
	err = retryer.RetryCallWithTimeoutErrorOnly(ctx, q.log, []int{1, 1, 3}, DatabaseErrorsToRetry, query)
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}

	return count, nil
}

func createDeletePersonStmt(ctx context.Context, tx *sql.Tx) (*sql.Stmt, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	psqlInsert, _, err := psql.Delete(GetTableFullName(PersonsTable)).
		Where(sq.Eq{"id": "?"}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("squirrel sql delete statement for '"+GetTableFullName(PersonsTable)+"': %w", err)

	}
	return tx.PrepareContext(ctx, psqlInsert)
}

func (q *QueriesHandler) UpdatePersonById(ctx context.Context, tx *sql.Tx, id int64, personData *datastructs.PersonData, genderId int64, countryId int64) (int64, error) {
	errorMsg := fmt.Sprintf("update person by id '%d' with data '%v' in '%s'", id, personData, PersonsTable) + ": %w"

	var columnsToUpdate []string
	for _, col := range ColumnsInPersonsTable {
		columnsToUpdate = append(columnsToUpdate, col)
	}

	stmt, err := createUpdatePersonByIdStmt(ctx, tx, columnsToUpdate)
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}
	defer helpers.ExecuteWithLogError(stmt.Close, q.log)

	var result sql.Result
	query := func(context context.Context) error {
		result, err = stmt.ExecContext(
			context,
			personData.Name,
			personData.Surname,
			personData.Patronymic,
			personData.Age,
			genderId,
			countryId,
			id,
		)
		return err
	}
	err = retryer.RetryCallWithTimeoutErrorOnly(ctx, q.log, []int{1, 1, 3}, DatabaseErrorsToRetry, query)
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}

	return count, nil
}

func (q *QueriesHandler) UpdatePartialPersonById(ctx context.Context, tx *sql.Tx, id int64, values map[string]interface{}) (int64, error) {
	errorMsg := fmt.Sprintf("update partially person by id '%d' with data '%v' in '%s'", id, values, PersonsTable) + ": %w"

	var columnsToUpdate []string
	var valuesToUpdate []interface{}
	for key, val := range values {
		columnsToUpdate = append(columnsToUpdate, key)
		valuesToUpdate = append(valuesToUpdate, val)
	}
	valuesToUpdate = append(valuesToUpdate, id)

	stmt, err := createUpdatePersonByIdStmt(ctx, tx, columnsToUpdate)
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}
	defer helpers.ExecuteWithLogError(stmt.Close, q.log)

	var result sql.Result
	query := func(context context.Context) error {
		result, err = stmt.ExecContext(
			context,
			valuesToUpdate...,
		)
		return err
	}
	err = retryer.RetryCallWithTimeoutErrorOnly(ctx, q.log, []int{1, 1, 3}, DatabaseErrorsToRetry, query)
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf(errorMsg, err)
	}

	return count, nil
}

func createUpdatePersonByIdStmt(ctx context.Context, tx *sql.Tx, values []string) (*sql.Stmt, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Update(GetTableFullName(PersonsTable))
	for _, col := range values {
		builder = builder.Set(col, "?")
	}
	builder = builder.Where(sq.Eq{"id": "?"})
	psqlUpdate, _, err := builder.ToSql()

	if err != nil {
		return nil, fmt.Errorf("squirrel sql update statement for '"+GetTableFullName(PersonsTable)+"': %w", err)

	}
	return tx.PrepareContext(ctx, psqlUpdate)
}

func (q *QueriesHandler) GetAdditionalId(ctx context.Context, tx *sql.Tx, name string, table string) (int64, error) {
	errorMsg := fmt.Sprintf("get additional id for '%s' in '%s'", name, table) + ": %w"

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
	errorMsg := fmt.Sprintf("insert additional value for '%s' in '%s'", name, table) + ": %w"

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
