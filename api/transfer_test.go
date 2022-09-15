package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	mock_db "github.com/faisal-a-n/simplebank/db/mock"
	db "github.com/faisal-a-n/simplebank/db/sqlc"
	"github.com/faisal-a-n/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateTxParamas struct {
	arg db.TransferTxParams
}

func (e eqCreateTxParamas) Matches(x interface{}) bool {
	arg, ok := x.(db.TransferTxParams)
	if !ok {
		return false
	}
	e.arg.CreatedAt = arg.CreatedAt

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateTxParamas) String() string {
	return fmt.Sprintf("matches arg  and response time")
}

func EqCreateUserParams(arg db.TransferTxParams) gomock.Matcher {
	return eqCreateTxParamas{arg}
}

func TestCreateTransferAPI(t *testing.T) {
	currency := "EUR"
	mismatchCurrency := "CAD"
	txAmount := int64(10)
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
				"amount":          txAmount,
				"currency":        currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, fromAccount.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(1).Return(toAccount, nil)

				fromEntry := db.Entry{
					ID:        1,
					AccountID: fromAccount.ID,
					Amount:    -txAmount,
					CreatedAt: time.Now().Unix(),
				}
				toEntry := db.Entry{
					ID:        2,
					AccountID: toAccount.ID,
					Amount:    txAmount,
					CreatedAt: time.Now().Unix(),
				}
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(1).Return(fromEntry, nil)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(1).Return(toEntry, nil)
				args := db.TransferTxParams{
					FromAccountID: fromAccount.ID,
					ToAccountID:   toAccount.ID,
					FromEntryID:   1,
					ToEntryID:     2,
					Amount:        txAmount,
					CreatedAt:     time.Now().Unix(),
				}
				store.EXPECT().TransferTx(gomock.Any(), EqCreateUserParams(args)).Times(1).Return(db.TransferTxResult{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "DifferentAccountOwner",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          txAmount,
				"currency":        currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				//Change the token id
				addAuthorizationHeader(t, request, maker, 3, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(1).Return(toAccount, nil)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "LowBalance",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          math.MaxInt64,
				"currency":        currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, fromAccount.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(1).Return(toAccount, nil)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "CurrencyMismatch",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          txAmount,
				"currency":        mismatchCurrency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, fromAccount.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(0)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "AccountNotFound",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          txAmount,
				"currency":        mismatchCurrency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, fromAccount.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(0)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: gin.H{
				"from_account_id": fromAccount.ID,
				"to_account_id":   toAccount.ID,
				"amount":          txAmount,
				"currency":        mismatchCurrency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, fromAccount.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(fromAccount.ID)).Times(1).Return(fromAccount, sql.ErrConnDone)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(toAccount.ID)).Times(0)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidBody",
			body: gin.H{
				"from_account_id": "fromAccount.ID",
				"to_account_id":   toAccount.ID,
				"amount":          txAmount,
				"currency":        mismatchCurrency,
			},
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorizationHeader(t, request, maker, fromAccount.UserID, authorizationHeaderKey, authorizationType, time.Minute)
			},
			builStubs: func(store *mock_db.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateEntry(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mock_db.NewMockStore(controller)
			testCase.builStubs(store)

			server := NewTestServer(t, store)
			url := "/transfers"
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
