package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/faisal-a-n/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func addAuthorizationHeader(t *testing.T, request *http.Request, maker token.Maker,
	userID int64, authKey string, authType string, durtation time.Duration) {
	access_token, err := maker.CreateToken(userID, durtation)
	require.NoError(t, err)
	require.NotEmpty(t, access_token)

	request.Header.Add(authKey, fmt.Sprintf("%s %s", authType, access_token))
}

func TestMiddleware(t *testing.T) {
	userID := int64(1)
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, userID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidHeader",
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, userID, authorizationHeaderKey, "", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthType",
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, userID, authorizationHeaderKey, "notBearer", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidToken",
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, userID, authorizationHeaderKey, "notBearer", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthHeader",
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			server := NewTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(authPath, authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}
