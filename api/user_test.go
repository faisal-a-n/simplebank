package api

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock_db "github.com/faisal-a-n/simplebank/db/mock"
	db "github.com/faisal-a-n/simplebank/db/sqlc"
	"github.com/faisal-a-n/simplebank/util"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {
	user := generateRandomUser()
	password := user.Password
	testCases := []struct {
		name          string
		body          db.User
		buildStub     func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: user,
			buildStub: func(store *mock_db.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "InvalidName",
			body: db.User{
				Name:     "#213csae",
				Password: password,
				Email:    user.Email,
			},
			buildStub: func(store *mock_db.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			body: db.User{
				Name:     user.Name,
				Password: "pass",
				Email:    user.Email,
			},
			buildStub: func(store *mock_db.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: db.User{
				Name:     user.Name,
				Password: password,
				Email:    "user.Email",
			},
			buildStub: func(store *mock_db.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "AccountWithEmailAlreadyRegistered",
			body: db.User{
				Name:     user.Name,
				Password: password,
				Email:    user.Email,
			},
			buildStub: func(store *mock_db.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			body: db.User{
				Name:     user.Name,
				Password: password,
				Email:    user.Email,
			},
			buildStub: func(store *mock_db.MockStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			store := mock_db.NewMockStore(mockController)
			testCase.buildStub(store)

			url := "/users"
			body, err := json.Marshal(map[string]string{
				"name":     testCase.body.Name,
				"email":    testCase.body.Email,
				"password": testCase.body.Password,
			})
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			server := NewServer(store)
			server.router.ServeHTTP(recorder, request)
		})
	}
}

func generateRandomUser() db.User {
	return db.User{
		ID:                util.GenerateRandomInt(1000, 1),
		Name:              util.GenerateString(8),
		Email:             util.RandomEmail(),
		Password:          util.GenerateString(8),
		PasswordChangedAt: time.Now().Unix(),
		CreatedAt:         time.Now().Unix(),
	}
}
