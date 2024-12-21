package gophermart

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/repository/mocks"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/pkg/jwt"
)

func TestGetUsersWithdrawHistory(t *testing.T) {
	type want struct {
		resultLen int
		err       error
	}

	type expects struct {
		repoError  error
		repoResult []repoModel.Withdraw
	}

	tests := []struct {
		name    string
		want    want
		expects expects
	}{
		{
			name: "successfully get users withdraw history",
			want: want{resultLen: 2},
			expects: expects{repoResult: []repoModel.Withdraw{
				{
					OrderID:  "1",
					Withdraw: decimal.NewFromFloat(10),
				},
				{
					OrderID:  "2",
					Withdraw: decimal.NewFromFloat(20),
				},
			}},
		},
		{
			name:    "get orders: repo error",
			want:    want{resultLen: 0},
			expects: expects{repoError: errors.New("repo error")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockGophermartRepository(ctrl)
			mock.EXPECT().GetUsersHistory(gomock.Any(), gomock.Any()).Return(tt.expects.repoResult, tt.expects.repoError)

			serv := NewGophermartService(mock, dummy.NewDummyLogger(), jwt.NewJWT(dummy.NewDummyLogger()))
			models, _ := serv.GetUsersWithdrawHistory(context.Background(), "1")

			assert.Equal(t, tt.want.resultLen, len(models))
		})
	}
}
