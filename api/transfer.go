package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	From_account_id int64  `json:"from_account_id" binding:"required,min=1"`
	To_account_id   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount          int64  `json:"amount" binding:"required,gt=0"`
	Currency        string `json:"currency" binding:"required,currency"`
}

// 创建转账记录,用shouldbindjson
func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//检查传输币种是否匹配
	fromAccount, valid := server.validAccount(ctx, req.From_account_id, req.Currency)
	if !valid {
		return
	}
	//要转账的账户不是这个登陆的账户
	authPayload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("要转账的账户不是这个登陆的账户")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.To_account_id, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTXParams{
		FromAccountID: req.From_account_id,
		ToAccountID:   req.To_account_id,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTX(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// 检查两个账户的currency是否相同
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	//err不为空存在两种情况，没有找到或者网络错误
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("账户[%d]币种不匹配[%s][%s]", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, false
}
