package service

import (
	"context"
	"fmt"

	"github.com/yukiHaga/go_todo_app/entity"
	"github.com/yukiHaga/go_todo_app/store"
)

type ListTasks struct {
	DB   store.Queryer
	Repo TaskLister
}

func (listTasks *ListTasks) ListTasks(ctx context.Context) (entity.Tasks, error) {
	ts, err := listTasks.Repo.ListTasks(ctx, listTasks.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return ts, nil
}
