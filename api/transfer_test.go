package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_db "github.com/faisal-a-n/simplebank/db/mock"
	"github.com/faisal-a-n/simplebank/token"
	"github.com/faisal-a-n/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferAPI(t *testing.T) {
	currency := util.GenerateCurrency()
	fromAccount := randomAccountWithCurrency(currency)
	toAccount := randomAccountWithCurrency(currency)
	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		builStubs     func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          10,
				"currency":        currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, fromAccount.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
			},
		},
	}
	for _, testCase := range testCases {
		break
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mock_db.NewMockStore(controller)

			server := NewTestServer(t, store)
			url := "/transfer"
			body, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)

			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}
