package handler

import (
	"context"

	"github.com/yukiHaga/go_todo_app/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService AddTaskService
type ListTasksService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}

// これはモックを自動生成するためのコメント
// go generateコマンドによって実行できる
// go性のツールならば、go installを使わずとも実行できる
// ビルドタグと一緒でgo generateの部分にスペース入れちゃダメ

// ハンドラーパッケージの中でサービスのインターフェースを定義するから、呼び出す時に、パッケージ名.をつけなくて済むから楽だ。だからか。
type RegisterUserService interface {
	RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error)
}

type LoginService interface {
	Login(ctx context.Context, name, pw string) (string, error)
}
