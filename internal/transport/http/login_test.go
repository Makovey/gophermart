package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Makovey/gophermart/internal/logger/dummy"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/mocks"
)

func TestLoginHandler(t *testing.T) {
	type want struct {
		code          int
		hasAuthHeader bool
	}

	type expects struct {
		expectServiceCall bool
		serviceError      error
	}

	type params struct {
		body io.Reader
	}

	tests := []struct {
		name    string
		want    want
		expects expects
		params  params
	}{
		{
			name: "success login",
			want: want{
				code: http.StatusOK,
			},
			expects: expects{
				expectServiceCall: true,
			},
			params: params{
				body: strings.NewReader(makeJSON(map[string]any{
					"login":    "testLogin",
					"password": "testPassword",
				})),
			},
		},
		{
			name: "error login: with body reader",
			want: want{
				code: http.StatusBadRequest,
			},
			expects: expects{},
			params: params{
				body: errReader(0),
			},
		},
		{
			name: "error login: with empty body",
			want: want{
				code: http.StatusBadRequest,
			},
			expects: expects{},
			params: params{
				body: strings.NewReader(makeJSON(map[string]any{})),
			},
		},
		{
			name: "error login: with long login or password",
			want: want{
				code: http.StatusBadRequest,
			},
			expects: expects{},
			params: params{
				body: strings.NewReader(makeJSON(map[string]any{
					"login":    strings.Repeat("l", 31),
					"password": strings.Repeat("p", 31),
				})),
			},
		},
		{
			name: "error register: user not found",
			want: want{
				code: http.StatusUnauthorized,
			},
			expects: expects{
				expectServiceCall: true,
				serviceError:      service.ErrNotFound,
			},
			params: params{
				body: strings.NewReader(makeJSON(map[string]any{
					"login":    "testLogin",
					"password": "testPassword",
				})),
			},
		},
		{
			name: "error register: password does not match",
			want: want{
				code: http.StatusUnauthorized,
			},
			expects: expects{
				expectServiceCall: true,
				serviceError:      service.ErrPasswordDoesntMatch,
			},
			params: params{
				body: strings.NewReader(makeJSON(map[string]any{
					"login":    "testLogin",
					"password": "testPassword",
				})),
			},
		},
		{
			name: "error login: service error",
			want: want{
				code: http.StatusInternalServerError,
			},
			expects: expects{
				expectServiceCall: true,
				serviceError:      service.ErrExecStmt,
			},
			params: params{
				body: strings.NewReader(makeJSON(map[string]any{
					"login":    "testLogin",
					"password": "testPassword",
				})),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv := mocks.NewMockGophermartService(ctrl)
			if tt.expects.expectServiceCall {
				serv.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Return(uuid.NewString(), tt.expects.serviceError)
			}

			h := NewHTTPHandler(
				dummy.NewDummyLogger(),
				serv,
			)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/user/login", tt.params.body)

			h.Login(w, r)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			if tt.expects.expectServiceCall && tt.expects.serviceError == nil {
				assert.Empty(t, resBody)
				assert.NotEmpty(t, w.Header().Get("Authorization"))
			} else {
				assert.Empty(t, w.Header().Get("Authorization"))
			}
		})
	}
}
