package service

import (
	"context"

	"github.com/yukiHaga/go_todo_app/entity"
	"github.com/yukiHaga/go_todo_app/store"
)

type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, task *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}

type UserRegister interface {
	RegisterUser(ctx context.Context, db store.Execer, u *entity.User) error
}

type UserGetter interface {
	GetUser(ctx context.Context, db store.Queryer, name string) (*entity.User, error)
}

type TokenGenerator interface {
	GenerateToken(ctx context.Context, u entity.User) ([]byte, error)
}
