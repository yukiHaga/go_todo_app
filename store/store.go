// 永続化の仮実装

package store

import (
	"errors"
	"sync"

	"github.com/yukiHaga/go_todo_app/entity"
)

var (
	Tasks       = &TaskStore{Tasks: map[entity.TaskID]*entity.Task{}}
	ErrNotFound = errors.New("not found")
)

type TaskStore struct {
	LastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
	Mutex  sync.Mutex
}

func (ts *TaskStore) Add(t *entity.Task) (*entity.Task, error) {
	ts.Mutex.Lock()

	ts.LastID++
	t.ID = ts.LastID
	ts.Tasks[t.ID] = t

	ts.Mutex.Unlock()
	return t, nil
}

// Allはタスク一覧を返す
func (ts *TaskStore) All() entity.Tasks {
	tasks := make([]*entity.Task, len(ts.Tasks))
	// タスクとスライスをマッピングさせる
	for _, t := range ts.Tasks {
		tasks[t.ID] = t
	}
	return tasks
}
