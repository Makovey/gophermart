package http

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/middleware"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/mocks"
	"github.com/golang/mock/gomock"
)

func TestPostOrderHandler(t *testing.T) {
	type want struct {
		code int
	}

	type expects struct {
		processCall  bool
		processErr   error
		validateCall bool
		isBodyValid  bool
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
			name: "success posting",
			want: want{
				code: http.StatusAccepted,
			},
			expects: expects{
				processCall:  true,
				validateCall: true,
				isBodyValid:  true,
			},
			params: params{
				body:      strings.NewReader("12345678903"),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting order: empty token",
			want: want{
				code: http.StatusBadRequest,
			},
			expects: expects{},
			params: params{
				body: strings.NewReader("12345678903"),
			},
		},
		{
			name: "error posting order: with body reader",
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
			name: "error posting order: with empty body",
			want: want{
				code: http.StatusBadRequest,
			},
			expects: expects{},
			params: params{
				body:      strings.NewReader(""),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting order: with invalid body",
			want: want{
				code: http.StatusUnprocessableEntity,
			},
			expects: expects{
				validateCall: true,
			},
			params: params{
				body:      strings.NewReader("1221"),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting order: process returned conflict error",
			want: want{
				code: http.StatusConflict,
			},
			expects: expects{
				processCall:  true,
				processErr:   service.ErrOrderConflict,
				validateCall: true,
				isBodyValid:  true,
			},
			params: params{
				body:      strings.NewReader("12345678903"),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting order: process returned already posted error",
			want: want{
				code: http.StatusOK,
			},
			expects: expects{
				processCall:  true,
				processErr:   service.ErrOrderAlreadyPosted,
				validateCall: true,
				isBodyValid:  true,
			},
			params: params{
				body:      strings.NewReader("12345678903"),
				authToken: uuid.NewString(),
			},
		},
		{
			name: "error posting order: process returned internal error",
			want: want{
				code: http.StatusInternalServerError,
			},
			expects: expects{
				processCall:  true,
				processErr:   service.ErrExecStmt,
				validateCall: true,
				isBodyValid:  true,
			},
			params: params{
				body:      strings.NewReader("12345678903"),
				authToken: uuid.NewString(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := mocks.NewMockOrderService(ctrl)
			if tt.expects.validateCall {
				serv.EXPECT().ValidateOrderID(gomock.Any()).Return(tt.expects.isBodyValid)
			}

			if tt.expects.processCall {
				serv.EXPECT().ProcessNewOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.expects.processErr)
			}

			h := NewHTTPHandler(
				dummy.NewDummyLogger(),
				mocks.NewMockUserService(ctrl),
				serv,
				mocks.NewMockBalanceService(ctrl),
				mocks.NewMockHistoryService(ctrl),
			)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/user/orders", tt.params.body)
			ctx := context.WithValue(r.Context(), middleware.CtxUserIDKey, tt.params.authToken)

			h.PostOrder(w, r.WithContext(ctx))

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
		})
	}
}
