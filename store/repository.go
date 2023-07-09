package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/yukiHaga/go_todo_app/clock"
	"github.com/yukiHaga/go_todo_app/config"
)

// DBへの接続情報を隠蔽する
// RDBを利用しなくなった、つまり、アプリケーションの終了のタイミングでDBのクローズをしたりは自動でできないので、
// Closeを呼び出す関数を返すようにしている
func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {

	// Openはデータベースドライバ名と、データソース名で指定されたデータベースを開く
	// 返されたDBは、複数のゴルーチンによる同時使用に対して安全であり、アイドル状態の接続プールを維持する
	// アイドル状態とは、何もしていないけど利用できる状態のこと
	// したがって、一度Open関数を呼べばよくて、DBを閉じる必要はない。
	// parseTime=trueを忘れると、time.Time型のフィールドに正しい時刻情報が取得できないので、注意
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	))
	if err != nil {
		return nil, nil, err
	}

	// Openは実際に接続テストが行われない
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// sql.Open関数は接続確認までは行わない。
	// そのため、明示的に*sql.DB.PingContextメソッドを利用して疎通確認できるかを確認している
	// PingContext はデータベースへの接続が生きているかどうかを確認し、必要であれば接続を確立します。
	if err := db.PingContext(ctx); err != nil {
		log.Println(err.Error())
		return nil, func() { db.Close() }, err
	}

	// NewDb は既存の *sql.DB の新しい sqlx DB ラッパーを返します。
	// 名前付きクエリのサポートには、元のデータベースの driverName が必要です。
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, func() { db.Close() }, nil
}

// Txは進行中のデータベーストランザクションである。
// トランザクションはコミットまたはロールバックの呼び出しで終了しなければならない。
type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Stmtは、sql.Stmtに追加機能を持たせたsqlxラッパーです。
type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

// 書き込み系の操作を集めたインターフェース
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// 参照系の操作を集めたインターフェース
// アプリケーションコードとしても「このメソッドの引数はQueryerインターフェースなので、MySQL上のデータを更新することはないな」と
// コードリーディングがしやすくなる
type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

var (
	// ある構造体が特定のインターフェースを満たしているかチェックしている
	// つまり、インターフェースが期待通りに定義されているかをチェックしている
	// 以下の場合、sqlx.DBがBegginerインターフェースを満たしているかチェックしている
	// nilを*sqlx.DB型でキャストしたものを特定のインターフェース型のブランク変数に代入しているのか。
	// *sqlx.DB型が正しくインターフェースが実装されていれば、互換性があるので代入できる
	// *sqlx.DBパッケージに既に実装されていた。
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.Tx)(nil)
)

// Clockerフィールドは、SQL実行時に利用する時刻情報を制御するためのclock.Clockerインターフェースである
// 永続化操作を行う際の、時刻を固定化できるようにするのが目的である。
type Repository struct {
	Clocker clock.Clocker
}

const (
	// ErrCodeMySQLDuplicateEntryはMySQL系のDUPLICATEエラーコード
	ErrCodeMySQLDuplicateEntry = 1062
)

var (
	ErrAlreadyEntry = errors.New("duplicate entry")
)
