package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/middleware"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service/mocks"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

func TestGetOrdersHandler(t *testing.T) {
	type want struct {
		code int
	}

	type expects struct {
		getOrdersCall bool
		orders        []model.OrderResponse
		ordersErr     error
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
				getOrdersCall: true,
				orders:        generateModels(2),
			},
			params: params{
				authToken: uuid.NewString(),
			},
		},
		{
			name: "success getting empty orders",
			want: want{
				code: http.StatusNoContent,
			},
			expects: expects{
				getOrdersCall: true,
				orders:        generateModels(0),
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
			expects: expects{
				orders: generateModels(0),
			},
			params: params{},
		},
		{
			name: "failed getting orders with service error",
			want: want{
				code: http.StatusInternalServerError,
			},
			expects: expects{
				getOrdersCall: true,
				ordersErr:     errors.New("service error"),
				orders:        generateModels(0),
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

			serv := mocks.NewMockOrderService(ctrl)

			if tt.expects.getOrdersCall {
				serv.EXPECT().GetOrders(gomock.Any(), gomock.Any()).Return(tt.expects.orders, tt.expects.ordersErr)
			}

			h := NewHTTPHandler(
				dummy.NewDummyLogger(),
				mocks.NewMockUserService(ctrl),
				serv,
				mocks.NewMockBalanceService(ctrl),
				mocks.NewMockHistoryService(ctrl),
			)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			ctx := context.WithValue(r.Context(), middleware.CtxUserIDKey, tt.params.authToken)

			h.GetOrders(w, r.WithContext(ctx))

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}

func generateModels(len int) []model.OrderResponse {
	var models []model.OrderResponse
	for i := 0; i < len; i++ {
		models = append(models, model.OrderResponse{
			Number:     strconv.Itoa(i),
			Status:     repoModel.New,
			Accrual:    nil,
			UploadedAt: time.Now().String(),
		})
	}

	return models
}
