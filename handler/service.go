package handler

import (
	"context"

	"github.com/yukiHaga/go_todo_app/entity"
)

// これはモックを自動生成するためのコメント
// go generateコマンドによって実行できる
// go性のツールならば、go installを使わずとも実行できる
// go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService AddTaskService
type ListTasksService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}
