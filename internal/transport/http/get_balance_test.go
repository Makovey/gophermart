package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/middleware"
	"github.com/Makovey/gophermart/internal/service/mocks"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

func TestGetBalanceHandler(t *testing.T) {
	type want struct {
		code int
	}

	type expects struct {
		getBalanceCall bool
		balanceRes     model.BalanceResponse
		balanceErr     error
	}

	type params struct {
		authToken string
	}

	tests := []struct {
		name    string
		want    want
		expects expects
		params  params
	}{
		{
			name: "successfully getting balance",
			want: want{
				code: http.StatusOK,
			},
			expects: expects{
				getBalanceCall: true,
			},
			params: params{
				authToken: uuid.NewString(),
			},
		},
		{
			name: "get balance error: empty token",
			want: want{
				code: http.StatusBadRequest,
			},
			expects: expects{},
			params:  params{},
		},
		{
			name: "get balance error: service error",
			want: want{
				code: http.StatusInternalServerError,
			},
			expects: expects{
				getBalanceCall: true,
				balanceErr:     errors.New("service error"),
			},
			params: params{
				authToken: uuid.NewString(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := mocks.NewMockGophermartService(ctrl)

			if tt.expects.getBalanceCall {
				serv.EXPECT().GetUsersBalance(gomock.Any(), gomock.Any()).Return(tt.expects.balanceRes, tt.expects.balanceErr)
			}

			h := NewHTTPHandler(
				dummy.NewDummyLogger(),
				serv,
			)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
			ctx := context.WithValue(r.Context(), middleware.CtxUserIDKey, tt.params.authToken)

			h.GetBalance(w, r.WithContext(ctx))

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
