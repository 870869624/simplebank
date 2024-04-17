package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount() //随机生成账户
	//测试用
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		//有该用户
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)). //具有任何上下文和次特定账户id的参数
					Times(1).                                        //执行一次
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code) //最简单的处理方式是检查相应码是否正确

				requireBdoymatchAccount(t, recoder.Body, account)
			},
		},
		//添加更多案例， 没有该用户
		{
			name:      "NotFounf",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)). //具有任何上下文和次特定账户id的参数
					Times(1).                                        //执行一次
					Return(db.Account{}, sql.ErrNoRows)              //空的账户结构体和没有查询到的错误返回
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recoder.Code) //最简单的处理方式是检查相应码是否正确
			},
		},
		//网络错误
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)). //具有任何上下文和次特定账户id的参数
					Times(1).                                        //执行一次
					Return(db.Account{}, sql.ErrConnDone)            //空的账户结构体和没有查询到的错误返回
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code) //最简单的处理方式是检查相应码是否正确
			},
		},
		//badrequest
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()). //具有任何上下文和次特定账户id的参数
					Times(0)                                //执行一次
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
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
			server := NewServer(store)

			recoder := httptest.NewRecorder() //返回一个初始化的回应请求记录器
			url := fmt.Sprintf("/accounts/%d", tc.accountID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request) //使用创建的记录和请求
			tc.checkResponse(t, recoder)
		})

	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInit(1, 100),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

// 检查响应体是否匹配
func requireBdoymatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account

	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func TestCreateAccount(t *testing.T) {

}

func TestListAccount(t *testing.T) {

}
