package handler

import (
	"net/http"
)

type ListTasks struct {
	// Store *store.TaskStore
	// DB   *sqlx.DB
	// Repo store.Repository
	Service ListTasksService
}

// func NewListTasks() *ListTasks {
// 	return &ListTasks{
// 		// Store: store.Tasks,
// 	}
// }

// SeerveHTTPメソッドを満たすことで、AddTaskがHandlerインターフェースを満たすことになる
func (listTasks *ListTasks) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// リストタスクハンドラーはDBとRepositoryに依存している必要があったけど、サービスだけに依存すれば良くなった
	// tasks, err := listTasks.Repo.ListTasks(ctx, listTasks.DB)
	tasks, err := listTasks.Service.ListTasks(ctx)
	if err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
	}
	RespondJson(ctx, w, tasks, http.StatusOK)
}
