package postgresql // import "github.com/erupshis/effective_mobile/internal/server/storage/managers/postgresql"

Package postgresql postgresql handling PostgreSQL database.

CONSTANTS

const (
	SchemaName     = "persons_data"
	PersonsTable   = "persons"
	GendersTable   = "genders"
	CountriesTable = "countries"
)

FUNCTIONS

func CreatePostgreDB(ctx context.Context, cfg config.Config, queriesHandler QueriesHandler, log logger.BaseLogger) (managers.BaseStorageManager, error)
    CreatePostgreDB creates manager implementation. Supports migrations and
    check connection to database.


TYPES

type QueriesHandler struct {
	// Has unexported fields.
}
    QueriesHandler support object that implements database's queries and
    responsible for connection to databse.

func CreateHandler(log logger.BaseLogger) QueriesHandler
    CreateHandler creates QueriesHandler.

func (q *QueriesHandler) DeletePerson(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
    DeletePerson performs direct query request to database to delete person by
    id.

func (q *QueriesHandler) GetAdditionalId(ctx context.Context, tx *sql.Tx, name string, table string) (int64, error)
    GetAdditionalId returns foreign key from linked table.

func (q *QueriesHandler) InsertAdditionalId(ctx context.Context, tx *sql.Tx, name string, table string) (int64, error)
    InsertAdditionalId adds new value forl linked table and returns foreign key
    from linked table.

func (q *QueriesHandler) InsertPerson(ctx context.Context, tx *sql.Tx, personData *datastructs.PersonData, genderId int64, countryId int64) (int64, error)
    InsertPerson performs direct query request to database to add new person.

func (q *QueriesHandler) SelectPersons(ctx context.Context, tx *sql.Tx, filters map[string]interface{}, pageNum int64, pageSize int64) ([]datastructs.PersonData, error)
    SelectPersons performs direct query request to database to select persons
    satisfying filters. Supports pagination.

func (q *QueriesHandler) UpdatePartialPersonById(ctx context.Context, tx *sql.Tx, id int64, values map[string]interface{}) (int64, error)
    UpdatePartialPersonById generates statement for delete query.

