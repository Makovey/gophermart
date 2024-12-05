package gophermart

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/repository/mocks"
	"github.com/Makovey/gophermart/internal/transport/http/model"
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
			name:    "Successful generate new authorization token",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "testablePassword"}},
			expects: expects{expectRepoCall: true},
		},
		{
			name:    "Failed generate new authorization token with repo error",
			param:   params{authModel: model.AuthRequest{Login: "testableLogin", Password: "testablePassword"}},
			expects: expects{expectRepoCall: true, repoError: errors.New("repoError")},
		},
		{
			name:    "Failed generate new authorization token with long password",
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

			serv := NewGophermartService(mock, dummy.NewDummyLogger())
			token, err := serv.RegisterUser(context.Background(), tt.param.authModel)

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
