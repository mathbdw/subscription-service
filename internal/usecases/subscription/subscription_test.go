package subscription

import (
	"context"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/google/uuid"
	"github.com/mathbdw/subscription-service/internal/domain/entities"
	"github.com/mathbdw/subscription-service/internal/errors"
	"github.com/mathbdw/subscription-service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	subTest = entities.Subscription{
		ID:          1,
		ServiceName: "Test service",
		UserId:      uuid.New(),
		Price:       100,
		StartDate:   time.Now(),
	}

	updateFields = map[string]any{
		"service_name": subTest.ServiceName,
		"user_id": subTest.UserId,
		"price": subTest.Price,
		"start_date": subTest.StartDate,
	}

	filterCost = entities.FilterParams{
		ServiceName: subTest.ServiceName,
		UserId: subTest.UserId,
		StartDate: entities.DateRange{From: &subTest.StartDate},
	}
)

func TestSubscription_Create_ErrorRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	mockSubRepo.EXPECT().
		Create(ctx, subTest).
		Return(errors.New("error repo"))

	err := us.Create(ctx, subTest)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "SubscriptionUsecase.Create: repo exec")
}

func TestSubscription_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	mockSubRepo.EXPECT().
		Create(ctx, subTest).
		Return(nil)

	err := us.Create(ctx, subTest)

	require.NoError(t, err)
}

func TestSubscription_GetByID_ErrorRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	mockSubRepo.EXPECT().
		GetByID(ctx, subTest.ID).
		Return(nil, errors.New("error repo"))

	sub, err := us.GetByID(ctx, subTest.ID)

	require.Error(t, err)
	require.Nil(t, sub)
	assert.Contains(t, err.Error(), "SubscriptionUsecase.GetByID: repo exec")
}

func TestSubscription_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	mockSubRepo.EXPECT().
		GetByID(ctx, subTest.ID).
		Return(&subTest, nil)

	sub, err := us.GetByID(ctx, subTest.ID)

	require.NoError(t, err)
	require.Equal(t, subTest.ID, sub.ID)
	require.Equal(t, subTest.ServiceName, sub.ServiceName)
}

func TestSubscription_List_ErrorRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	params := entities.QueryCriteria{}

	mockSubRepo.EXPECT().
		List(ctx, gomock.Any()).
		Return(nil, errors.New("error repo"))

	sub, err := us.List(ctx, params)

	require.Error(t, err)
	require.Nil(t, sub)
	assert.Contains(t, err.Error(), "SubscriptionUsecase.List: repo exec")
}

func TestSubscription_List_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	params := entities.QueryCriteria{}
	respData := &entities.ResponseListSubscription{}

	mockSubRepo.EXPECT().
		List(ctx, gomock.Any()).
		Return(respData, nil)

	res, err := us.List(ctx, params)

	require.NoError(t, err)
	require.Equal(t, respData, res)
}

func TestSubscription_Update_ErrorRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	mockSubRepo.EXPECT().
		GetByID(ctx, subTest.ID).
		Return(&subTest, nil)
		
	mockSubRepo.EXPECT().
		Update(ctx, subTest.ID, updateFields).
		Return(errors.New("error repo"))

	err := us.Update(ctx, subTest.ID, updateFields)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "SubscriptionUsecase.Update: repo exec")
}

func TestSubscription_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	mockSubRepo.EXPECT().
		GetByID(ctx, subTest.ID).
		Return(&subTest, nil)

	mockSubRepo.EXPECT().
		Update(ctx, subTest.ID, updateFields).
		Return(nil)

	err := us.Update(ctx, subTest.ID, updateFields)

	require.NoError(t, err)
}

func TestSubscription_Delete_ErrorRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	mockSubRepo.EXPECT().
		GetByID(ctx, subTest.ID).
		Return(&subTest, nil)

	mockSubRepo.EXPECT().
		Delete(ctx, subTest.ID).
		Return(errors.New("error repo"))

	err := us.Delete(ctx, subTest.ID)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "SubscriptionUsecase.Delete: repo exec")
}

func TestSubscription_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	mockSubRepo.EXPECT().
		GetByID(ctx, subTest.ID).
		Return(&subTest, nil)

	mockSubRepo.EXPECT().
		Delete(ctx, subTest.ID).
		Return(nil)

	err := us.Delete(ctx, subTest.ID)

	require.NoError(t, err)
}

func TestSubscription_GetCost_ErrorRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	mockSubRepo.EXPECT().
		GetCost(ctx, filterCost).
		Return(int64(0), errors.New("error repo"))

	resCost, err := us.GetCost(ctx, filterCost)

	require.Error(t, err)
	require.Equal(t, int64(0), resCost)
	assert.Contains(t, err.Error(), "SubscriptionUsecase.GetCost: repo exec")
}

func TestSubscription_GetCost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSubRepo := mocks.NewMockSubscriptionRepository(ctrl)
	mockLogger := mocks.NewMockLogger(ctrl)
	us := NewSubscriptionUsecase(mockSubRepo, mockLogger)
	ctx := context.Background()

	mockSubRepo.EXPECT().
		GetCost(ctx, filterCost).
		Return(int64(124), nil)

	resCost, err := us.GetCost(ctx, filterCost)

	require.NoError(t, err)
	require.Equal(t, int64(124), resCost)
}
