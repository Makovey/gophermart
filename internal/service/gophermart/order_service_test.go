package gophermart

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/repository/mocks"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/pkg/jwt"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
)

func TestValidateOrderID(t *testing.T) {
	type params struct {
		orderID string
	}

	type expects struct {
		expectedValid bool
	}

	tests := []struct {
		name    string
		expects expects
		params  params
	}{
		{
			name:    "successful validate order id",
			expects: expects{expectedValid: true},
			params:  params{orderID: "12345678903"},
		},
		{
			name:    "error validate with zero",
			expects: expects{expectedValid: false},
			params:  params{orderID: "0"},
		},
		{
			name:    "error validate with negative number",
			expects: expects{expectedValid: false},
			params:  params{orderID: "-12345678903"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockGophermartRepository(ctrl)

			serv := NewGophermartService(mock, dummy.NewDummyLogger(), jwt.NewJWT(dummy.NewDummyLogger()))
			isValid := serv.ValidateOrderID(tt.params.orderID)

			if tt.expects.expectedValid {
				assert.True(t, isValid)
			} else {
				assert.False(t, isValid)
			}
		})
	}
}

func TestProcessNewOrder(t *testing.T) {
	type want struct {
		finalErr error
	}

	type params struct {
		userID  string
		orderID string
	}

	type expects struct {
		getError    error
		getCallAns  repoModel.Order
		postNewCall bool
		postError   error
	}

	tests := []struct {
		name    string
		want    want
		param   params
		expects expects
	}{
		{
			name:    "process order: posted new order",
			want:    want{finalErr: nil},
			param:   params{userID: "12345", orderID: "1"},
			expects: expects{getError: service.ErrNotFound, postNewCall: true, postError: nil},
		},
		{
			name:    "process order: posted new order with error",
			want:    want{finalErr: service.ErrExecStmt},
			param:   params{userID: "12345", orderID: "1"},
			expects: expects{getError: service.ErrNotFound, postNewCall: true, postError: service.ErrExecStmt},
		},
		{
			name:    "process error: order already posted by another user",
			want:    want{finalErr: service.ErrOrderConflict},
			param:   params{userID: "12345", orderID: "1"},
			expects: expects{getCallAns: repoModel.Order{OwnerUserID: "1"}},
		},
		{
			name:    "process order: already posted by user",
			want:    want{finalErr: service.ErrOrderAlreadyPosted},
			param:   params{userID: "12345", orderID: "1"},
			expects: expects{getCallAns: repoModel.Order{OwnerUserID: "12345"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockGophermartRepository(ctrl)
			mock.EXPECT().GetOrderByID(gomock.Any(), tt.param.orderID).Return(tt.expects.getCallAns, tt.expects.getError)

			if tt.expects.postNewCall {
				mock.EXPECT().PostNewOrder(gomock.Any(), tt.param.orderID, tt.param.userID).Return(tt.expects.postError)
			}

			serv := NewGophermartService(mock, dummy.NewDummyLogger(), jwt.NewJWT(dummy.NewDummyLogger()))
			err := serv.ProcessNewOrder(context.Background(), tt.param.userID, tt.param.orderID)

			if tt.want.finalErr != nil {
				assert.Equal(t, tt.want.finalErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetOrders(t *testing.T) {
	type want struct {
		resultLen int
		err       bool
	}

	type params struct {
		userID string
	}

	type expects struct {
		repoError  error
		repoResult []repoModel.Order
	}

	tests := []struct {
		name    string
		want    want
		param   params
		expects expects
	}{
		{
			name:    "successful get orders",
			want:    want{resultLen: 5},
			param:   params{userID: "12345"},
			expects: expects{repoResult: generateModels(5)},
		},
		{
			name:    "get orders: repo error",
			want:    want{resultLen: 0},
			param:   params{userID: "12345"},
			expects: expects{repoError: errors.New("err")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockGophermartRepository(ctrl)
			mock.EXPECT().GetOrders(gomock.Any(), tt.param.userID).Return(tt.expects.repoResult, tt.expects.repoError)

			serv := NewGophermartService(mock, dummy.NewDummyLogger(), jwt.NewJWT(dummy.NewDummyLogger()))
			models, err := serv.GetOrders(context.Background(), tt.param.userID)

			if tt.want.err {
				assert.Error(t, err)
				assert.Nil(t, models)
			}

			assert.Equal(t, tt.want.resultLen, len(models))
		})
	}
}

func generateModels(len int) []repoModel.Order {
	var models []repoModel.Order
	for i := 0; i < len; i++ {
		models = append(models, repoModel.Order{
			OrderID:     strconv.Itoa(i),
			OwnerUserID: strconv.Itoa(i),
			Status:      repoModel.New,
			Accrual:     newPointerDec(i),
			CreatedAt:   time.Now(),
		})
	}

	return models
}

func newPointerDec(val int) *decimal.Decimal {
	d := new(decimal.Decimal)
	i := decimal.NewFromInt(int64(val))
	d = &i
	return d
}
