package gophermart

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/repository/mocks"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	servMock "github.com/Makovey/gophermart/internal/service/mocks"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/pkg/jwt"
)

func TestRegisterNewUser(t *testing.T) {
	type params struct {
		authModel model.AuthRequest
	}

	type expects struct {
		repoCall  bool
		repoError error
	}

	tests := []struct {
		name    string
		param   params
		expects expects
	}{
		{
			name:    "successful generate new authorization token",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "testablePassword"}},
			expects: expects{repoCall: true},
		},
		{
			name:    "failed generate new authorization token with repo error",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "testablePassword"}},
			expects: expects{repoCall: true, repoError: errors.New("repoError")},
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

			mock := mocks.NewMockUserServiceRepository(ctrl)
			if tt.expects.repoCall {
				mock.EXPECT().RegisterNewUser(gomock.Any(), gomock.Any()).Return(tt.expects.repoError)
			}

			serv := NewGophermartService(
				NewUserService(mock, jwt.NewJWT(dummy.NewDummyLogger())),
				servMock.NewMockOrderService(ctrl),
				servMock.NewMockBalanceService(ctrl),
				servMock.NewMockHistoryService(ctrl),
			)
			token, err := serv.RegisterNewUser(context.Background(), tt.param.authModel)

			if err != nil {
				assert.Empty(t, token)
			} else {
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
			expects: expects{repoAnswer: repoModel.RegisterUser{UserID: "id", Login: "testableLogin", PasswordHash: "testablePassword"}},
		},
		{
			name:    "failed login: repo error",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "testablePassword"}},
			expects: expects{repoError: errors.New("repoError")},
		},
		{
			name:    "failed  login: password does not match",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "newPassword"}},
			expects: expects{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			pass, _ := bcrypt.GenerateFromPassword([]byte(tt.expects.repoAnswer.PasswordHash), bcrypt.DefaultCost)
			tt.expects.repoAnswer.PasswordHash = string(pass)

			mock := mocks.NewMockUserServiceRepository(ctrl)
			mock.EXPECT().LoginUser(gomock.Any(), tt.param.authModel.Login).Return(tt.expects.repoAnswer, tt.expects.repoError)

			serv := NewGophermartService(
				NewUserService(mock, jwt.NewJWT(dummy.NewDummyLogger())),
				servMock.NewMockOrderService(ctrl),
				servMock.NewMockBalanceService(ctrl),
				servMock.NewMockHistoryService(ctrl),
			)
			token, err := serv.LoginUser(context.Background(), tt.param.authModel)

			if err != nil {
				assert.Empty(t, token)
			} else {
				assert.NotEmpty(t, token)
			}
		})
	}
}
