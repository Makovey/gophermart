package gophermart

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/repository/mocks"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/pkg/jwt"
)

func TestGeneratePasswordHash(t *testing.T) {
	type params struct {
		authModel model.AuthRequest
	}

	type expects struct {
		expectRepoCall bool
		repoError      error
	}

	tests := []struct {
		name    string
		param   params
		expects expects
	}{
		{
			name:    "successful generate new authorization token",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "testablePassword"}},
			expects: expects{expectRepoCall: true},
		},
		{
			name:    "failed generate new authorization token with repo error",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "testablePassword"}},
			expects: expects{expectRepoCall: true, repoError: errors.New("repoError")},
		},
		{
			name:    "failed generate new authorization token with long password",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: strings.Repeat("Password", 10)}},
			expects: expects{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockGophermartRepository(ctrl)
			if tt.expects.expectRepoCall {
				mock.EXPECT().RegisterNewUser(gomock.Any(), gomock.Any()).Return(tt.expects.repoError)
			}

			serv := NewGophermartService(mock, dummy.NewDummyLogger(), jwt.NewJWT(dummy.NewDummyLogger()))
			token, err := serv.RegisterNewUser(context.Background(), tt.param.authModel)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	type params struct {
		authModel model.AuthRequest
	}

	type expects struct {
		repoCall   bool
		repoError  error
		repoAnswer repoModel.RegisterUser
	}

	tests := []struct {
		name    string
		param   params
		expects expects
	}{
		{
			name:    "successful generate new authorization token",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "testablePassword"}},
			expects: expects{repoCall: true, repoAnswer: repoModel.RegisterUser{UserID: "id", Login: "testableLogin", PasswordHash: "testablePassword"}},
		},
		{
			name:    "failed login: repo error",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "testablePassword"}},
			expects: expects{repoCall: true, repoError: errors.New("repoError")},
		},
		{
			name:    "failed  login: password does not match",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "newPassword"}},
			expects: expects{repoCall: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pass, _ := bcrypt.GenerateFromPassword([]byte(tt.expects.repoAnswer.PasswordHash), bcrypt.DefaultCost)
			tt.expects.repoAnswer.PasswordHash = string(pass)

			mock := mocks.NewMockGophermartRepository(ctrl)
			if tt.expects.repoCall {
				mock.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(tt.expects.repoAnswer, tt.expects.repoError)
			}

			serv := NewGophermartService(mock, dummy.NewDummyLogger(), jwt.NewJWT(dummy.NewDummyLogger()))
			token, err := serv.LoginUser(context.Background(), tt.param.authModel)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

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
		getCall     bool
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
			expects: expects{getCall: true, getError: service.ErrNotFound, postNewCall: true, postError: nil},
		},
		{
			name:    "process order: posted new order with error",
			want:    want{finalErr: service.ErrExecStmt},
			param:   params{userID: "12345", orderID: "1"},
			expects: expects{getCall: true, getError: service.ErrNotFound, postNewCall: true, postError: service.ErrExecStmt},
		},
		{
			name:    "process error: order already posted by another user",
			want:    want{finalErr: service.ErrOrderConflict},
			param:   params{userID: "12345", orderID: "1"},
			expects: expects{getCall: true, getCallAns: repoModel.Order{OwnerUserID: "1"}},
		},
		{
			name:    "process order: already posted by user",
			want:    want{finalErr: service.ErrOrderAlreadyPosted},
			param:   params{userID: "12345", orderID: "1"},
			expects: expects{getCall: true, getCallAns: repoModel.Order{OwnerUserID: "12345"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockGophermartRepository(ctrl)
			if tt.expects.getCall {
				mock.EXPECT().GetOrderByID(gomock.Any(), gomock.Any()).Return(tt.expects.getCallAns, tt.expects.getError)
			}
			if tt.expects.postNewCall {
				mock.EXPECT().PostNewOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expects.postError)
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
