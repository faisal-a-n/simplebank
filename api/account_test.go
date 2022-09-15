package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_db "github.com/faisal-a-n/simplebank/db/mock"
	db "github.com/faisal-a-n/simplebank/db/sqlc"
	"github.com/faisal-a-n/simplebank/token"
	"github.com/faisal-a-n/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		builStubs     func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Ok",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, account.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkAccounts(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, account.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, account.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, account.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		//TODO: Add more cases
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()
			store := mock_db.NewMockStore(mockController)
			testCase.builStubs(store)

			//Make http request
			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", testCase.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, req, server.tokenMaker)

			server.router.ServeHTTP(recorder, req)

			testCase.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()
	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		builStubs     func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"name":     account.Name,
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, account.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				args := db.CreateAccountParams{
					UserID:    account.UserID,
					Name:      account.Name,
					Currency:  account.Currency,
					Balance:   0,
					CreatedAt: account.CreatedAt,
				}
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				fmt.Println(recorder.Body)
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "InvalidRequestBody",
			body: gin.H{
				"name":     account.Name,
				"currency": "currency",
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, account.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"name":     account.Name,
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, account.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				args := db.CreateAccountParams{
					UserID:    account.UserID,
					Name:      account.Name,
					Currency:  account.Currency,
					CreatedAt: account.CreatedAt,
				}
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "AccountWithCurrencyAlreadyExists",
			body: gin.H{
				"name":     account.Name,
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, account.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "UserError",
			body: gin.H{
				"name":     account.Name,
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, account.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, &pq.Error{Code: "23503"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			store := mock_db.NewMockStore(mockController)
			testCase.builStubs(store)
			//Make http request
			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/accounts"
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)
			testCase.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)

			testCase.checkResponse(t, recorder)
		})
	}
}

func TestGetAccountsAPI(t *testing.T) {
	n := 1
	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount()
	}
	userID := accounts[len(accounts)-1].UserID

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		buildStub     func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"count":   5,
				"page_id": 1,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, userID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			buildStub: func(store *mock_db.MockStore) {
				args := db.ListAccountsForUserParams{
					UserID: userID,
					Limit:  5,
					Offset: 0,
				}
				store.EXPECT().ListAccountsForUser(gomock.Any(), gomock.Eq(args)).Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response map[string]json.RawMessage
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)

				var fetchedAccounts []db.Account
				err = json.Unmarshal(response["data"], &fetchedAccounts)
				require.NoError(t, err)
				require.NotEmpty(t, fetchedAccounts)

				for _, v := range fetchedAccounts {
					require.Equal(t, userID, v.UserID)
				}
			},
		},
		{
			name: "InvalidQuery",
			body: gin.H{
				"count":   3,
				"page_id": 0,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, userID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			buildStub: func(store *mock_db.MockStore) {
				store.EXPECT().ListAccountsForUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"count":   5,
				"page_id": 1,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, userID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			buildStub: func(store *mock_db.MockStore) {
				store.EXPECT().ListAccountsForUser(gomock.Any(), gomock.Any()).Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "NoAccountsAvailable",
			body: gin.H{
				"count":   10,
				"page_id": 5,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, userID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			buildStub: func(store *mock_db.MockStore) {
				store.EXPECT().ListAccountsForUser(gomock.Any(), gomock.Any()).Times(1).
					Return([]db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			store := mock_db.NewMockStore(mockController)

			testCase.buildStub(store)

			url := fmt.Sprintf("/accounts?page_id=%d&count=%d", testCase.body["page_id"], testCase.body["count"])

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			server := NewTestServer(t, store)
			require.NoError(t, err)
			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}

}

func randomAccount() db.Account {
	user := generateRandomUser()
	return db.Account{
		ID:        util.GenerateRandomInt(1000, 1),
		Name:      util.GenerateString(6),
		Balance:   util.GenerateAmount(),
		Currency:  util.GenerateCurrency(),
		UserID:    user.ID,
		CreatedAt: time.Now().Unix(),
	}
}

func randomAccountWithCurrency(currency string) (account db.Account) {
	account = randomAccount()
	account.Currency = currency
	return
}

func checkAccounts(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var response map[string]json.RawMessage
	err = json.Unmarshal(data, &response)
	require.NoError(t, err)

	var fetchedAccount db.Account
	err = json.Unmarshal(response["data"], &fetchedAccount)
	require.NoError(t, err)

	require.Equal(t, account, fetchedAccount)
}
