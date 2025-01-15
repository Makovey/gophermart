package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/middleware"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/mocks"
)

func TestPostWithdrawHandler(t *testing.T) {
	type want struct {
		code int
	}

	type expects struct {
		serviceCall bool
		serviceErr  error
	}

	type params struct {
		body      io.Reader
		authToken string
	}

	tests := []struct {
		name    string
		want    want
		expects expects
		params  params
	}{
		{
			name: "successful post withdraw",
			want: want{
				code: http.StatusOK,
			},
			expects: expects{
				serviceCall: true,
			},
			params: params{
				body: strings.NewReader(makeJSON(map[string]any{
					"order": "1234567890",
					"sum":   100.50,
				})),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting withdraw: empty token",
			want: want{
				code: http.StatusBadRequest,
			},
			expects: expects{},
			params:  params{},
		},
		{
			name: "error posting withdraw: with body reader",
			want: want{
				code: http.StatusBadRequest,
			},
			expects: expects{},
			params: params{
				body:      errReader(0),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting withdraw: without order id",
			want: want{
				code: http.StatusBadRequest,
			},
			expects: expects{},
			params: params{
				body: strings.NewReader(makeJSON(map[string]any{
					"sum": 100.50,
				})),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting withdraw: with invalid body",
			want: want{
				code: http.StatusInternalServerError,
			},
			expects: expects{},
			params: params{
				body:      strings.NewReader("1221"),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting withdraw: service return random error",
			want: want{
				code: http.StatusInternalServerError,
			},
			expects: expects{
				serviceCall: true,
				serviceErr:  errors.New("random error"),
			},
			params: params{

				body: strings.NewReader(makeJSON(map[string]any{
					"order": "1234567890",
					"sum":   100.50,
				})),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting withdraw: service return not enough founds error",
			want: want{
				code: http.StatusPaymentRequired,
			},
			expects: expects{
				serviceCall: true,
				serviceErr:  service.ErrNotEnoughFounds,
			},
			params: params{

				body: strings.NewReader(makeJSON(map[string]any{
					"order": "1234567890",
					"sum":   100.50,
				})),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting withdraw: service return not found order",
			want: want{
				code: http.StatusUnprocessableEntity,
			},
			expects: expects{
				serviceCall: true,
				serviceErr:  service.ErrNotFound,
			},
			params: params{

				body: strings.NewReader(makeJSON(map[string]any{
					"order": "1234567890",
					"sum":   100.50,
				})),
				authToken: uuid.NewString(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := mocks.NewMockBalanceService(ctrl)
			if tt.expects.serviceCall {
				serv.EXPECT().WithdrawUsersBalance(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expects.serviceErr)
			}

			h := NewHTTPHandler(
				dummy.NewDummyLogger(),
				mocks.NewMockUserService(ctrl),
				mocks.NewMockOrderService(ctrl),
				serv,
				mocks.NewMockHistoryService(ctrl),
			)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", tt.params.body)
			ctx := context.WithValue(r.Context(), middleware.CtxUserIDKey, tt.params.authToken)

			h.PostWithdraw(w, r.WithContext(ctx))

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
