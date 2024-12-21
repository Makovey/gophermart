package gophermart

import (
	"context"
	"errors"
	"github.com/Makovey/gophermart/internal/service"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/repository/mocks"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/internal/types"
	"github.com/Makovey/gophermart/pkg/jwt"
	"github.com/golang/mock/gomock"
)

func TestBalanceServiceGetUsersBalance(t *testing.T) {
	type want struct {
		balance model.BalanceResponse
	}

	type expects struct {
		repoError  error
		repoResult repoModel.Balance
	}

	tests := []struct {
		name    string
		want    want
		expects expects
	}{
		{
			name: "successfully fetched balance",
			want: want{model.BalanceResponse{
				Current:   types.FloatDecimal(decimal.NewFromFloat(20)),
				Withdrawn: types.FloatDecimal(decimal.NewFromFloat(5)),
			}},
			expects: expects{
				repoResult: repoModel.Balance{
					Accrual:   decimal.NewFromFloat(20),
					Withdrawn: decimal.NewFromFloat(5),
				}},
		},
		{
			name: "get users balance error: repo error",
			want: want{model.BalanceResponse{}},
			expects: expects{
				repoError:  errors.New("repo error"),
				repoResult: repoModel.Balance{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockGophermartRepository(ctrl)
			mock.EXPECT().GetUsersBalance(gomock.Any(), gomock.Any()).Return(tt.expects.repoResult, tt.expects.repoError)

			serv := NewGophermartService(mock, dummy.NewDummyLogger(), jwt.NewJWT(dummy.NewDummyLogger()))
			mod, _ := serv.GetUsersBalance(context.Background(), "1")

			assert.Equal(t, mod, tt.want.balance)
		})
	}
}

func TestWithdrawUsersBalance(t *testing.T) {
	type want struct {
		serviceError error
	}

	type expects struct {
		getUsersBalance    bool
		decreaseBalance    bool
		recordUsersBalance bool
	}

	type repoResult struct {
		getOrderErr        error
		getOrderResponse   repoModel.Order
		getBalanceErr      error
		getBalanceResult   repoModel.Balance
		decreaseBalanceErr error
		recordHistoryErr   error
	}

	type args struct {
		userID  string
		request model.WithdrawRequest
	}

	tests := []struct {
		name       string
		want       want
		expects    expects
		repoResult repoResult
		args       args
	}{
		{
			name: "successfully withdraw balance",
			want: want{},
			expects: expects{
				getUsersBalance:    true,
				decreaseBalance:    true,
				recordUsersBalance: true,
			},
			repoResult: repoResult{getOrderResponse: repoModel.Order{OwnerUserID: "1"}},
			args:       args{userID: "1"},
		},
		{
			name:       "withdraw error: order conflict",
			want:       want{serviceError: service.ErrOrderConflict},
			expects:    expects{},
			repoResult: repoResult{getOrderResponse: repoModel.Order{OwnerUserID: "2"}},
			args:       args{userID: "1"},
		},
		{
			name:       "withdraw error: balance repo error",
			want:       want{serviceError: service.ErrNotFound},
			expects:    expects{getUsersBalance: true},
			repoResult: repoResult{getBalanceErr: service.ErrNotFound},
			args:       args{},
		},
		{
			name:       "withdraw error: accrual on balance less than withdraw",
			want:       want{serviceError: service.ErrNotEnoughFounds},
			expects:    expects{getUsersBalance: true},
			repoResult: repoResult{getBalanceResult: repoModel.Balance{Accrual: decimal.NewFromFloat(100)}},
			args:       args{request: model.WithdrawRequest{Sum: decimal.NewFromFloat(150)}},
		},
		{
			name:       "withdraw error: decrease balance repo error",
			want:       want{serviceError: service.ErrNotFound},
			expects:    expects{getUsersBalance: true, decreaseBalance: true},
			repoResult: repoResult{decreaseBalanceErr: service.ErrNotFound},
			args:       args{},
		},
		{
			name: "withdraw error: record withdraw history repo error",
			want: want{serviceError: service.ErrNotFound},
			expects: expects{
				getUsersBalance:    true,
				decreaseBalance:    true,
				recordUsersBalance: true,
			},
			repoResult: repoResult{recordHistoryErr: service.ErrNotFound},
			args:       args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockGophermartRepository(ctrl)
			mock.EXPECT().GetOrderByID(gomock.Any(), gomock.Any()).Return(tt.repoResult.getOrderResponse, tt.repoResult.getOrderErr)
			if tt.expects.getUsersBalance {
				mock.EXPECT().GetUsersBalance(gomock.Any(), gomock.Any()).Return(tt.repoResult.getBalanceResult, tt.repoResult.getBalanceErr)
			}

			if tt.expects.decreaseBalance {
				mock.EXPECT().DecreaseUsersBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.repoResult.decreaseBalanceErr)
			}

			if tt.expects.recordUsersBalance {
				mock.EXPECT().RecordUsersWithdraw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.repoResult.recordHistoryErr)
			}

			serv := NewGophermartService(mock, dummy.NewDummyLogger(), jwt.NewJWT(dummy.NewDummyLogger()))
			err := serv.WithdrawUsersBalance(context.Background(), tt.args.userID, tt.args.request)

			assert.Equal(t, err, tt.want.serviceError)
		})
	}
}
