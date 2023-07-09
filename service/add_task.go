// service.AddTask型もstoreパッケージの特定の型に依存せず、インターフェースをDIする設計になっている
package service

import (
	"context"
	"fmt"

	"github.com/yukiHaga/go_todo_app/entity"
	"github.com/yukiHaga/go_todo_app/store"
)

// Execerはインターフェース
// サービスがstore.Repositoryに直接依存しないようにインターフェースを定義した
// Repository.AddTaskをこのパッケージの中で直接使っちゃうと、Repository.AddTaskの戻り値や引数が変化したりしたら、常にサービスは影響を受ける。
// 共通のインターフェースを作っておいて、それにserviceとstore.Repositoryが従うことで、serviceはstore.Repostiroyへの影響を気にしなくて良くなる
// 逆も言えて、store.Repositoryがserivceの変更を気にしなくて良くなる。そのため、インターフェースを挟むことで他のモジュールを気にしなくても、自分のモジュールを変更できる

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

func (addTask *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	// タスクの初期化をサービスに移した
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	// タスクの登録
	err := addTask.Repo.AddTask(ctx, addTask.DB, t)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return t, nil
}
