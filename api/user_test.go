package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/harrychopra/go-api/db/mock"
	db "github.com/harrychopra/go-api/db/models"
	"github.com/harrychopra/go-api/util"
	"github.com/stretchr/testify/require"
)

// Customer matcher implementation to test CreateUser Handler
// eqCreateUserParamsMatcher is our "wanted"/ "expected" object
type eqCreateUserParamsMatcher struct {
	expUser     db.CreateUserParams
	expPassword string
}

// Matches matches expected CreateUser() parameters with an actual argument
func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	// x is the "actual" sent by handler
	gotUser, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	// Verify that the "expected" naked password string, matches the hashed password
	if err := util.CheckPassword(e.expPassword, gotUser.HashedPassword); err != nil {
		return false
	}
	// Populate the empty HashedPassword field in our "expected" type from the "actual" object
	e.expUser.HashedPassword = gotUser.HashedPassword
	return reflect.DeepEqual(e.expUser, gotUser)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches expected arg %v and password: %v", e.expUser, e.expPassword)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{
		expUser:     arg,
		expPassword: password,
	}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser()
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mock.MockStore) {
				// Expected object
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mock.NewMockStore(ctrl)
			testCase.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := "/users"
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}
}

func randomUser() (db.User, string) {
	return db.User{
		Username: util.RandomName(),
		FullName: util.RandomName(),
		Email:    util.RandomEmail(),
	}, util.RandomString(8)
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)
}
