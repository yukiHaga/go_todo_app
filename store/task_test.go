package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/yukiHaga/go_todo_app/clock"
	"github.com/yukiHaga/go_todo_app/entity"
	"github.com/yukiHaga/go_todo_app/testutil"
)

// 実際のRDBMSを使ったテスト
func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()

	// entity.Taskを作成する他のテストケースと混ざるとテストがフェイルする
	// そのため、トランザクションを張ることで、このテストケースの中だけのテーブル状態にする
	// BeginTxx はトランザクションを開始し、*sql.Tx の代わりに *sqlx.Tx を返す。
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)

	// このテストケースが完了したら元に戻す
	// ロールバックはトランザクションを中止する
	// ロールバックなので、おそらくDBの状態もトランザクションの開始直前の状態に戻っている
	t.Cleanup(func() { tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	wants := prepareTasks(ctx, t, tx)

	// ListTasksを使うためにリポジトリーを定義した
	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

// モックを使ったテスト、DBへは直接接続せず、SQLドライバの結果をシミュレートしている
// そのおかげで、テストのDBクセス時間を減らせる
func TestRepository_AddTask(t *testing.T) {
	// Parallel は、このテストが他の並列テストと（そして他の並列テストとだけ）並列に実行されることを示します。
	// (厳密には並行だけど)
	// データベースへアクセスするようなテストは、データベースのアクセス２時間がかかるので、そこを並行でやれるとテスト時間を短縮できる
	// あるパッケージ内のテストは逐次的に実行されるけど、あるパッケージと別のパッケージのテストは並行に実施される
	// パッケージ内のテストでテスト関数ごとに並列を実現したいなら、t.Parallelを使う
	// t.Parallelを書いたテスト関数は実行時に一時停止して、パッケージ内のテストが全て完了した際に、再開して並行に実行される
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	task := &entity.Task{
		Title:     "ok task",
		Status:    "todo",
		CreatedAt: c.Now(),
		UpdatedAt: c.Now(),
	}

	// 空のモックを作成
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	// go-sqlmockを使うことで、実際のデータベース接続を必要とせずに、テスト内でSQLドライバの動作をシミュレートすることできる。
	// ExpectExec関数は、期待するUPDATE、INSERT、DELETE文と返り値の組み合わせをモック化します。
	// ExpectQuery関数ほど、返り値は重要ではありません。
	// ExpectExec関数の引数には、期待するSQLクエリ(UPDATE、INSERT、DELETE文)を指定します
	// ExpectExecはExec()がexpectedSQL queryで呼び出されることを期待します。
	// ExpectExec は、データベースの応答をモックすることができます。
	mock.ExpectExec(
		// エスケープが必要
		// これがexpected SQL query
		`INSERT INTO tasks \(title, status, created_at, updated_at\) VALUES \(\?, \?, \?, \?\)`,
	// WithArgsは与えられた期待される引数と実際のデータベース実行操作の引数をマッチさせます。
	).WithArgs(task.Title, task.Status, task.CreatedAt, task.UpdatedAt).
		// SQLmock.NewResult(lastInsertID int64, affectedRows int64)メソッドがあり、対応する結果を作成します。
		// NewResult は、Exec ベースのクエリ・モック用の新しい sql ドライバ結果を作成します。
		// NewResult関数の引数には、期待するカラムの値を指定します
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, task); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

// タスクテーブルをリセットして、テストデータを入れるための関数
func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()

	// 一度キレイにしておく
	if _, err := con.ExecContext(ctx, "DELETE FROM tasks;"); err != nil {
		t.Logf("failed to initialize tasks: %v", err)
	}

	// テスト用の固定された時刻を生成するクロッカーを生成
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title:     "want task 1",
			Status:    "todo",
			CreatedAt: c.Now(),
			UpdatedAt: c.Now(),
		},
		{
			Title:     "want task 2",
			Status:    "todo",
			CreatedAt: c.Now(),
			UpdatedAt: c.Now(),
		},
		{
			Title:     "want task 3",
			Status:    "done",
			CreatedAt: c.Now(),
			UpdatedAt: c.Now(),
		},
	}

	sql := `INSERT INTO tasks (title, status, created_at, updated_at)
	        VALUES
			(?, ?, ?, ?),
			(?, ?, ?, ?),
			(?, ?, ?, ?);`
	result, err := con.ExecContext(ctx, sql,
		wants[0].Title, wants[0].Status, wants[0].CreatedAt, wants[0].UpdatedAt,
		wants[1].Title, wants[1].Status, wants[1].CreatedAt, wants[1].UpdatedAt,
		wants[2].Title, wants[2].Status, wants[2].CreatedAt, wants[2].UpdatedAt,
	)
	if err != nil {
		t.Fatal(err)
	}

	// 複数のレコードをインサートで作成した時のsql.Result.LastInsertIdメソッドの戻り値となるIDは、
	// MySQLでは一つ目のレコードのID(発行されたIDの中で一番小さいID)になることが注意点です
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}
