package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// 自定义匹配器
type eqCreateUserparamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

// 检测输入参数与预制参数是否匹配
func (e eqCreateUserparamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserparamsMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%s)", e.arg, e.password)
}

// 定义的一个接口，返回的interface包含Matches
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserparamsMatcher{arg: arg, password: password}
}
func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t) //随机生成账户
	//测试用例，表示网络反映
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore) //模拟接口
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		//有该用户
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{ //arg为预期值
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)). //具有任何上下文和次特定账户id的参数
					Times(1).                                                    //执行一次
					Return(user, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code) //最简单的处理方式是检查相应码是否正确
				requireBdoymatchUser(t, recoder.Body, user)
			},
		},

		//网络错误
		{
			name: "InternalError",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()). //具有任何上下文和次特定账户id的参数
					Times(1).                               //执行一次
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code) //最简单的处理方式是检查相应码是否正确
			},
		},
		//重复用户名
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()). //具有任何上下文和次特定账户id的参数
					Times(1).                               //执行一次
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recoder.Code) //最简单的处理方式是检查相应码是否正确
			},
		},
		//不合法的用户名
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "Invalid_user#1",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()). //具有任何上下文和次特定账户id的参数
					Times(0)                                //执行一次
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code) //最简单的处理方式是检查相应码是否正确
			},
		},

		{
			name: "Invalidpassword",
			body: gin.H{
				"username":  user.Username,
				"password":  "11231",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()). //具有任何上下文和次特定账户id的参数
					Times(0)                                //执行一次
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code) //最简单的处理方式是检查相应码是否正确
			},
		},

		{
			name: "InvalidEmail",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "Invalid_Email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()). //具有任何上下文和次特定账户id的参数
					Times(0)                                //执行一次
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code) //最简单的处理方式是检查相应码是否正确
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t) //生成控制器
			defer ctrl.Finish()             //检查所有的操作是否已经完成/

			store := mockdb.NewMockStore(ctrl) //新建一个存储器，存储各项功能

			//构建存根
			tc.buildStubs(store)
			//启动服务发送请求
			server := NewTestServer(t, store)

			recoder := httptest.NewRecorder() //返回一个初始化的回应请求记录器

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			url := "/users"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request) //使用创建的记录和请求
			tc.checkResponse(recoder)
		})

	}

}

// 随机生成用户
func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashPassword,
		Email:          util.RandomEmail(),
		FullName:       util.RandomOwner(),
	}
	return
}

func requireBdoymatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotuser db.User
	err = json.Unmarshal(data, &gotuser)
	require.NoError(t, err)

	require.Equal(t, user.FullName, gotuser.FullName)
	require.Equal(t, gotuser.Username, user.Username)
	require.Empty(t, gotuser.HashedPassword)

	require.Equal(t, gotuser.Email, user.Email)

}
