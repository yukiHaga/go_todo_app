package store

import (
	"context"

	"github.com/yukiHaga/go_todo_app/entity"
)

// 全てのタスクを取得するメソッド
// 参照系のメソッドのため、Queryerインターフェースを満たす型の値を一つ受け取る
// Repositoryパターンを使うことで、永続化処理を抽象化している
// Repositoryを経由しないとデータにアクセスでない
// Repository層を経由してデータアクセスすることで、このメソッドを使う人は、どんなDBを使用しているのか、どんなSQL文なのか、どんなORMを使用しているのかを
// 知らなくてもよくなる。てか隠蔽できるので変更に強くなる。
// dbは引数で渡さなくても良い気もするけどな。dbの型はインターフェースだから何かしら理由があるのかも。
func (r *Repository) ListTasks(ctx context.Context, db Queryer) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT id, title, status, created_at, updated_at
	        FROM tasks;`
	// SelectContextメソッドは、sqlxパッケージの拡張メソッド
	// SelectContextメソッドは、複数のレコードを取得して、各レコードを一つ一つの構造体に代入したスライスを返してくれる
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}

// タスクを保存するメソッド
// タスクの保存に失敗したらエラーを返す
// タスクを戻り値で返しても良いのかなとも思ったが、railsのupdateだと真偽値しか返していないので、それを参照した
func (r *Repository) AddTask(ctx context.Context, db Execer, task *entity.Task) error {
	task.CreatedAt = r.Clocker.Now()
	task.UpdatedAt = r.Clocker.Now()

	sql := `INSERT INTO tasks
	        (title, status, created_at, updated_at)
			VALUES (:title, :status, :created_at, :updated_at);`
	result, err := db.NamedExecContext(ctx, sql, task)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// IDフィールドを更新することで、呼び出し元で定義したタスクに発行されたIDを伝える
	task.ID = entity.TaskID(id)

	return nil
}
