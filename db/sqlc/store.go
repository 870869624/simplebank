package db

import (
	"context"
	"database/sql"
	"fmt"
)

// 为了实现模拟数据库创建的借口，实现了查询的所有方法，也实现了交易的方法
type Store interface {
	TransferTX(ctx context.Context, arg TransferTXParams) (TransferTXResult, error)
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// 在数据库事务中执行函数
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	//如果错误不为零，就回滚事物
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v. rb err:%v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// 包含输入参数和转移交易
type TransferTXParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// 转移交易的结果
type TransferTXResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txkey = struct{}{}

// 具体的支付数据库处理事件， 创建了支付记录，更新账户金额
func (store *SQLStore) TransferTX(ctx context.Context, arg TransferTXParams) (TransferTXResult, error) {
	var result TransferTXResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txkey)
		fmt.Println(txName, "创建交易")

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "创建条目1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "创建条目2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    +arg.Amount,
		})
		if err != nil {
			return err
		}

		// result.FromAccount, err = q.AddAccountBalance()
		//TODO uodata account balance
		//从数据库获取余额，然后再增加或者减少一些金额，然后再更新回数据库

		if arg.FromAccountID < arg.ToAccountID {
			//存在潜在的死锁情况，两个函数更新同一个账户时
			fmt.Println(txName, "更新账户1余额")
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

		} else {
			fmt.Println(txName, "更新账户2余额")
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})

	return result, err
}

// 金额变更模块
func AddMoney(ctx context.Context, q *Queries, acountID1 int64, amount1 int64, acountID2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     acountID1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     acountID2,
	})
	if err != nil {
		return
	}
	return
}
