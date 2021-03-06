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

	"github.com/golang/mock/gomock"
	"github.com/harrychopra/go-api/db/mock"
	db "github.com/harrychopra/go-api/db/models"
	"github.com/harrychopra/go-api/token"
	"github.com/harrychopra/go-api/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	user, _ := randomUser()
	account := randomAccount(user.Username)
	// Table driven test set
	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenmaker token.Maker)
		buildStubs    func(store *mock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock.MockStore) {
				// Expecting api.GetAccount to call (once) the mocked Store.GetAccount() with a specific id
				// Mocked store.GetAccount returns an account object and nil error
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Verify the http response
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
			},
			// Middleware will abort this request and will not reach the Handler
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "Not Found",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock.MockStore) {
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
			name:      "Internal Error",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone) // Should be considered as an internal error
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Invalid ID_Bad Request",
			accountID: 0, // Account handler should invalidate the request as min accepted ID value is 1
			setupAuth: func(t *testing.T, request *http.Request, tokenmaker token.Maker) {
				addAuthorization(t, request, tokenmaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mock.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0) // Since ID is invalid, mocked GetAccount should not be called by Handler
				// Also No return value expected
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mock.NewMockStore(ctrl)

			// Build stubs for each test case
			testCase.buildStubs(store)

			// Start http server and send a test request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", testCase.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)

			// Send request to the server.router to serve
			server.router.ServeHTTP(recorder, request)

			// Verify the http response
			testCase.checkResponse(t, recorder)
		})
	}
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Currency: util.RandomCurrency(),
		Balance:  util.RandomMoney(),
	}
}

// requireBodyMatch asserts the response's body object against an input account object
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
