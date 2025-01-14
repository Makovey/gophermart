package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/middleware"
	"github.com/Makovey/gophermart/internal/service/mocks"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/internal/types"
)

func TestGetWithdrawsHistoryHandler(t *testing.T) {
	type want struct {
		code int
	}

	type expects struct {
		withdrawServiceCall bool
		serviceResp         []model.WithdrawHistoryResponse
		serviceErr          error
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
			name: "success getting orders",
			want: want{
				code: http.StatusOK,
			},
			expects: expects{
				withdrawServiceCall: true,
				serviceResp: []model.WithdrawHistoryResponse{
					{
						Order:       "1234",
						Sum:         types.FloatDecimal(decimal.NewFromFloat(10)),
						ProcessedAt: time.Now(),
					},
					{
						Order:       "4321",
						Sum:         types.FloatDecimal(decimal.NewFromFloat(20)),
						ProcessedAt: time.Now(),
					},
				},
			},
			params: params{
				authToken: uuid.NewString(),
			},
		},
		{
			name: "success getting empty history",
			want: want{
				code: http.StatusNoContent,
			},
			expects: expects{
				withdrawServiceCall: true,
			},
			params: params{
				authToken: uuid.NewString(),
			},
		},
		{
			name: "failed getting orders with empty token",
			want: want{
				code: http.StatusBadRequest,
			},
			expects: expects{},
			params:  params{},
		},
		{
			name: "withdraw history err: service error",
			want: want{
				code: http.StatusInternalServerError,
			},
			expects: expects{
				withdrawServiceCall: true,
				serviceErr:          errors.New("service error"),
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

			if tt.expects.withdrawServiceCall {
				serv.EXPECT().GetUsersWithdrawHistory(gomock.Any(), gomock.Any()).Return(tt.expects.serviceResp, tt.expects.serviceErr)
			}

			h := NewHTTPHandler(
				dummy.NewDummyLogger(),
				serv,
			)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
			ctx := context.WithValue(r.Context(), middleware.CtxUserIDKey, tt.params.authToken)

			h.GetWithdrawsHistory(w, r.WithContext(ctx))

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
