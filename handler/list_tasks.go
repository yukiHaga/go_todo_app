package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/yukiHaga/go_todo_app/store"
)

type ListTasks struct {
	// Store *store.TaskStore
	DB   *sqlx.DB
	Repo store.Repository
}

// func NewListTasks() *ListTasks {
// 	return &ListTasks{
// 		// Store: store.Tasks,
// 	}
// }

// SeerveHTTPメソッドを満たすことで、AddTaskがHandlerインターフェースを満たすことになる
func (listTasks *ListTasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := listTasks.Repo.ListTasks(ctx, listTasks.DB)
	if err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
	}
	RespondJson(ctx, w, tasks, http.StatusOK)
}
