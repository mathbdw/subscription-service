package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/mathbdw/subscription-service/internal/domain/entities"
	"github.com/mathbdw/subscription-service/internal/errors"
	"github.com/mathbdw/subscription-service/mocks"
)

var subTest = entities.Subscription{
	ServiceName: "test service",
	UserID:      uuid.New(),
	Price:       100,
	StartDate:   time.Now(),
}

var defConfigTest = defaultConfig

func setVar() {
	defaultConfig = defConfigTest
}

// Modifies the global variable table, so it cannot be executed in parallel.
// nolint:paralleltest
func TestUser_Create_ErrorBuilder(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO (price,service_name,start_date,user_id) VALUES ($1,$2,$3,$4)")).
		WithArgs(subTest.Price, subTest.ServiceName, subTest.StartDate, subTest.UserID).
		WillReturnError(errors.New("build query"))

	defaultConfig.table = ""
	err = repo.Create(ctx, entities.Subscription{
		ServiceName: subTest.ServiceName,
		UserID:      subTest.UserID,
		Price:       subTest.Price,
		StartDate:   subTest.StartDate,
		EndDate:     subTest.EndDate,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Create: build query")
}

// Modifies the global variable columnsSelect, so it cannot be executed in parallel.
// nolint:paralleltest
func TestUser_GetByID_ErrorBuilder(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	defaultConfig.columnsSelect = []string{}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT FROM WHERE id = $1)")).
		WithArgs(subTest.ID).
		WillReturnError(errors.New("build query"))

	sub, err := repo.GetByID(ctx, subTest.ID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.getByID: build query")
	assert.Nil(t, sub)
}

// Modifies the global variable columnsSelectCount, so it cannot be executed in parallel.
// nolint:paralleltest
func TestUser_List_ErrorBuildQueryCount(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM ")).
		WithoutArgs().
		WillReturnError(errors.New("build query"))

	defaultConfig.columnsSelectCount = []string{}
	respSubs, err := repo.List(ctx, entities.QueryCriteria{})

	require.Error(t, err)
	require.Nil(t, respSubs)
	assert.Contains(t, err.Error(), "subscriptionRepositories.List: build query count()")
}

// Modifies the global variable columnsSelect, so it cannot be executed in parallel.
// nolint:paralleltest
func TestUser_List_ErrorBuildQuery(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	qc := entities.QueryCriteria{
		Pagination: entities.PaginationParams{Page: uint64(2), Limit: uint64(2)},
	}

	limit := qc.Pagination.Limit
	offset := (qc.Pagination.Page - 1) * limit

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM subscription")).
		WithoutArgs().
		WillReturnRows(mock.NewRows([]string{"COUNT(*)"}).AddRow(uint64(10)))

	setVar()
	defaultConfig.columnsSelect = []string{}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT  FROM subscription LIMIT $1 OFFSET $2")).
		WithArgs(offset, limit).
		WillReturnError(errors.New("build query"))

	respSubs, err := repo.List(ctx, qc)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.List: build query")
	assert.Nil(t, respSubs)
}

// Modifies the global variable table, so it cannot be executed in parallel.
// nolint:paralleltest
func TestUser_Update_ErrorBuildQuery(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()
	fieldsUpdate := map[string]any{
		"service_name": "Test service",
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE  SET service_name = $1, updated_at = $2 WHERE id = $3")).
		WithArgs(fieldsUpdate["service_name"], sqlmock.AnyArg(), subTest.ID).
		WillReturnError(sql.ErrNoRows)

	defaultConfig.table = ""
	err = repo.Update(ctx, subTest.ID, fieldsUpdate)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Update: build query")
}

// Modifies the global variable table, so it cannot be executed in parallel.
// nolint:paralleltest
func TestUser_Delete_ErrorBuildQuery(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM  WHERE id = $1")).
		WithArgs(subTest.ID).
		WillReturnError(sql.ErrNoRows)

	defaultConfig.table = ""
	err = repo.Delete(ctx, subTest.ID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Delete: build query")
}

// Modifies the global variable columnsCost, so it cannot be executed in parallel.
// nolint:paralleltest
func TestUser_GetCost_ErrorBuildQuery(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT  FROM subscription")).
		WithoutArgs().
		WillReturnError(errors.New("build query"))

	defaultConfig.columnsCost = []string{}
	cost, err := repo.GetCost(ctx, entities.FilterParams{})

	require.Error(t, err)
	require.Equal(t, int64(0), cost)
	assert.Contains(t, err.Error(), "subscriptionRepositories.GetCost: build query")

	setVar()
}

func TestUser_Create_ErrorExec(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO subscription (price,service_name,start_date,user_id) VALUES ($1,$2,$3,$4)")).
		WithArgs(subTest.Price, subTest.ServiceName, subTest.StartDate, subTest.UserID).
		WillReturnError(sql.ErrNoRows)

	err = repo.Create(ctx, entities.Subscription{
		ServiceName: subTest.ServiceName,
		UserID:      subTest.UserID,
		Price:       subTest.Price,
		StartDate:   subTest.StartDate,
		EndDate:     subTest.EndDate,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Create: exec query")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Create custom ErrorResult.
type ErrorResult struct{}

func (r *ErrorResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (r *ErrorResult) RowsAffected() (int64, error) {
	return 0, errors.New("rows affected error")
}

func TestUser_Create_ErrorAffectedRows(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO subscription (price,service_name,start_date,user_id) VALUES ($1,$2,$3,$4)")).
		WithArgs(subTest.Price, subTest.ServiceName, subTest.StartDate, subTest.UserID).
		WillReturnResult(&ErrorResult{})

	err = repo.Create(ctx, entities.Subscription{
		ServiceName: subTest.ServiceName,
		UserID:      subTest.UserID,
		Price:       subTest.Price,
		StartDate:   subTest.StartDate,
		EndDate:     subTest.EndDate,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Create: get affected rows")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_Create_ErrorNotEquilRowsAffected(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO subscription (price,service_name,start_date,user_id) VALUES ($1,$2,$3,$4)")).
		WithArgs(subTest.Price, subTest.ServiceName, subTest.StartDate, subTest.UserID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Create(ctx, entities.Subscription{
		ServiceName: subTest.ServiceName,
		UserID:      subTest.UserID,
		Price:       subTest.Price,
		StartDate:   subTest.StartDate,
		EndDate:     subTest.EndDate,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Create: expected rowsAffected")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_Create_Success(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO subscription (price,service_name,start_date,user_id) VALUES ($1,$2,$3,$4)")).
		WithArgs(subTest.Price, subTest.ServiceName, subTest.StartDate, subTest.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Create(ctx, entities.Subscription{
		ServiceName: subTest.ServiceName,
		UserID:      subTest.UserID,
		Price:       subTest.Price,
		StartDate:   subTest.StartDate,
		// EndDate:     subTest.endTime,
	})

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_GetByID_ErrorScan(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, service_name, user_id, price, start_date, end_date FROM subscription WHERE id = $1")).
		WithArgs(subTest.ID).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetByID(ctx, subTest.ID)

	require.NoError(t, mock.ExpectationsWereMet())
	require.Error(t, err)
	assert.Contains(t, err.Error(), errors.ErrNotFound.Error())
	assert.Empty(t, user)
}

func TestUser_GetByID_Success(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, service_name, user_id, price, start_date, end_date FROM subscription WHERE id = $1")).
		WithArgs(subTest.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "service_name", "user_id", "price", "start_date", "end_date"}).
				AddRow(subTest.ID, subTest.ServiceName, subTest.UserID, subTest.Price, subTest.StartDate, subTest.EndDate),
		)

	model, err := repo.GetByID(ctx, subTest.ID)

	require.NoError(t, mock.ExpectationsWereMet())
	assert.NoError(t, err)
	assert.Equal(t, subTest.ID, model.ID)
	assert.Equal(t, subTest.ServiceName, model.ServiceName)
	assert.Equal(t, subTest.UserID, model.UserID)
	assert.Equal(t, subTest.Price, model.Price)
	assert.Equal(t, subTest.StartDate, model.StartDate)
	assert.Equal(t, subTest.EndDate, model.EndDate)
}

func TestUser_List_ErrorScanQueryCount(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM subscription")).
		WithoutArgs().
		WillReturnError(sql.ErrNoRows)

	respSubs, err := repo.List(ctx, entities.QueryCriteria{})

	require.Error(t, err)
	require.Nil(t, respSubs)
	assert.Contains(t, err.Error(), "subscriptionRepositories.List: scan query")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_List_ErrorGetQuery(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	qc := entities.QueryCriteria{
		Pagination: entities.PaginationParams{Page: uint64(2), Limit: uint64(2)},
	}

	limit := qc.Pagination.Limit
	offset := (qc.Pagination.Page - 1) * limit

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM subscription")).
		WithoutArgs().
		WillReturnRows(mock.NewRows([]string{"COUNT(*)"}).AddRow(uint64(10)))

	// defaultConfig.columnsSelect = []string{"id", "service_name", "user_id", "price", "start_date", "end_date"}
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT id, service_name, user_id, price, start_date, end_date FROM subscription LIMIT %d OFFSET %d", limit, offset))).
		WithoutArgs().
		WillReturnError(sql.ErrNoRows)

	respSubs, err := repo.List(ctx, qc)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.List: get query")
	assert.Nil(t, respSubs)
}

func TestUser_List_ErrorScanQuery(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	qc := entities.QueryCriteria{
		Pagination: entities.PaginationParams{Page: uint64(2), Limit: uint64(2)},
	}

	limit := qc.Pagination.Limit
	offset := (qc.Pagination.Page - 1) * limit

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM subscription")).
		WithoutArgs().
		WillReturnRows(mock.NewRows([]string{"COUNT(*)"}).AddRow(uint64(10)))

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT id, service_name, user_id, price, start_date, end_date FROM subscription LIMIT %d OFFSET %d", limit, offset))).
		WithoutArgs().
		WillReturnRows(mock.NewRows([]string{"id", "service_name", "user_id", "price", "start_date", "end_date"}).
			AddRow("", subTest.ServiceName, subTest.UserID, subTest.Price, subTest.StartDate, subTest.EndDate),
		)

	respSubs, err := repo.List(ctx, qc)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.List: scan query")
	assert.Nil(t, respSubs)
}

func TestUser_List_ErrorIterationQuery(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	qc := entities.QueryCriteria{
		Pagination: entities.PaginationParams{Page: uint64(2), Limit: uint64(2)},
	}

	limit := qc.Pagination.Limit
	offset := (qc.Pagination.Page - 1) * limit

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM subscription")).
		WithoutArgs().
		WillReturnRows(mock.NewRows([]string{"COUNT(*)"}).AddRow(uint64(10)))

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT id, service_name, user_id, price, start_date, end_date FROM subscription LIMIT %d OFFSET %d", limit, offset))).
		WithoutArgs().
		WillReturnRows(mock.NewRows([]string{"id", "service_name", "user_id", "price", "start_date", "end_date"}).
			AddRow(subTest.ID, subTest.ServiceName, subTest.UserID, subTest.Price, subTest.StartDate, subTest.EndDate).
			RowError(0, errors.New("network error")),
		)

	respSubs, err := repo.List(ctx, qc)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.List: iteration rows")
	assert.Nil(t, respSubs)
}

func TestUser_List_Success(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	qc := entities.QueryCriteria{
		Pagination: entities.PaginationParams{Page: uint64(2), Limit: uint64(2)},
	}

	limit := qc.Pagination.Limit
	offset := (qc.Pagination.Page - 1) * limit

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM subscription")).
		WithoutArgs().
		WillReturnRows(mock.NewRows([]string{"COUNT(*)"}).AddRow(uint64(3)))

	limit--
	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT id, service_name, user_id, price, start_date, end_date FROM subscription LIMIT %d OFFSET %d", limit, offset))).
		WithoutArgs().
		WillReturnRows(mock.NewRows([]string{"id", "service_name", "user_id", "price", "start_date", "end_date"}).
			AddRow(subTest.ID, subTest.ServiceName, subTest.UserID, subTest.Price, subTest.StartDate, subTest.EndDate),
		)

	respSubs, err := repo.List(ctx, qc)

	require.NoError(t, err)
	require.Len(t, respSubs.Data, 1)
}

func TestUser_Update_ValidateFalse(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	logger.EXPECT().Error(gomock.Any(), gomock.Any())
	fieldsUpdate := map[string]any{
		"service_name": "Test service",
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE subscription SET service_name = $1, updated_at = $2 WHERE id = $3")).
		WithArgs(fieldsUpdate["service_name"], sqlmock.AnyArg(), subTest.ID).
		WillReturnError(sql.ErrNoRows)

	fieldsUpdateIncorrect := map[string]any{
		"service_name": uint16(1000),
	}

	err = repo.Update(ctx, subTest.ID, fieldsUpdateIncorrect)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Update: validate")
}

func TestUser_Update_ErrorExecQuery(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	fieldsUpdate := map[string]any{
		"service_name": "Test service",
	}
	mock.ExpectExec(regexp.QuoteMeta("UPDATE subscription SET service_name = $1, updated_at = $2 WHERE id = $3")).
		WithArgs(fieldsUpdate["service_name"], sqlmock.AnyArg(), subTest.ID).
		WillReturnError(sql.ErrNoRows)

	err = repo.Update(ctx, subTest.ID, fieldsUpdate)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Update: exec query")
}

func TestUser_Update_ErrorAffectedRows(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	fieldsUpdate := map[string]any{
		"service_name": "Test service",
	}
	mock.ExpectExec(regexp.QuoteMeta("UPDATE subscription SET service_name = $1, updated_at = $2 WHERE id = $3")).
		WithArgs(fieldsUpdate["service_name"], sqlmock.AnyArg(), subTest.ID).
		WillReturnResult(&ErrorResult{})

	err = repo.Update(ctx, subTest.ID, fieldsUpdate)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Update: get affected rows")
}

func TestUser_Update_ErrorNotEquilRowsAffected(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	fieldsUpdate := map[string]any{
		"service_name": "Test service",
	}
	mock.ExpectExec(regexp.QuoteMeta("UPDATE subscription SET service_name = $1, updated_at = $2 WHERE id = $3")).
		WithArgs(fieldsUpdate["service_name"], sqlmock.AnyArg(), subTest.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Update(ctx, subTest.ID, fieldsUpdate)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Update: expected rowsAffected")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_Update_Success(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	fieldsUpdate := map[string]any{
		"service_name": "Test service",
	}
	mock.ExpectExec(regexp.QuoteMeta("UPDATE subscription SET service_name = $1, updated_at = $2 WHERE id = $3")).
		WithArgs(fieldsUpdate["service_name"], sqlmock.AnyArg(), subTest.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(ctx, subTest.ID, fieldsUpdate)

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_Delete_ErrorExecQuery(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM subscription WHERE id = $1")).
		WithArgs(subTest.ID).
		WillReturnError(sql.ErrNoRows)

	err = repo.Delete(ctx, subTest.ID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Delete: exec query")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_Delete_ErrorAffectedRows(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM subscription WHERE id = $1")).
		WithArgs(subTest.ID).
		WillReturnResult(&ErrorResult{})

	err = repo.Delete(ctx, subTest.ID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Delete: get affected rows")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_Delete_ErrorNotEquilRowsAffected(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM subscription WHERE id = $1")).
		WithArgs(subTest.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(ctx, subTest.ID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "subscriptionRepositories.Delete: expected rowsAffected")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_Delete_Success(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM subscription WHERE id = $1")).
		WithArgs(subTest.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(ctx, subTest.ID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_GetCost_ErrorScan(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	params := entities.FilterParams{}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(price), 0) FROM subscription")).
		WithoutArgs().
		WillReturnError(sql.ErrNoRows)

	cost, err := repo.GetCost(ctx, params)

	require.Error(t, err)
	require.Equal(t, int64(0), cost)
	assert.Contains(t, err.Error(), "subscriptionRepositories.GetCost: scan query")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUser_GetCost_Success(t *testing.T) {
	t.Parallel()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err, "create mock")
	defer func() {
		_ = mockDB.Close()
	}()

	ctrl := gomock.NewController(t)
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	logger := mocks.NewMockLogger(ctrl)
	repo := NewSubscriptionRepository(sqlxDB, builder, logger)
	ctx := context.Background()

	params := entities.FilterParams{}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(price), 0) FROM subscription")).
		WithoutArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"COALESCE(SUM(price)"}).
				AddRow(int64(124)),
		)

	cost, err := repo.GetCost(ctx, params)

	require.NoError(t, err)
	require.Equal(t, int64(124), cost)
	assert.NoError(t, mock.ExpectationsWereMet())
}
