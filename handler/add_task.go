package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AddTask struct {
	// DB        *sqlx.DB
	// Repo      store.Repository
	Service   AddTaskService
	Validator *validator.Validate
}

// DBを指定しちゃうと、モックがやりづらくなるな(それって、ハンドラーの部分だけDBを指定可能にすれば良いってことなのかな？)
// ただ、mysqlがいろんなところに書かれているのもいややな。変更に弱いし、MySQLじゃなくなった時に変更料が多すぎる。隠蔽したい。モックの時だけ利用できるようにしたい
// storeのNew関数を使えば、一応そういうことが起きないようにはなっているのか
// NewAddTaskを使うことで、フィールドが大幅に変更したとしても、既存のフィールドをなるべく変更せずに変更箇所を一箇所に集中できる
// func NewAddTask(db store.Execer) *AddTask {
// 	return &AddTask{
// 		// Store:     store.Tasks,
// 		DB:        d,
// 		Repo:      &store.Repository{},
// 		Validator: validator.New(),
// 	}
// }

// SeerveHTTPメソッドを満たすことで、AddTaskがHandlerインターフェースを満たすことになる
func (addTask *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// リクエストパラメータ
	// リクエストにバリデーションをかけたいので、リクエスストパラメータをバリデーションで定義する
	// バリデーション対象の構造体にvalidateタグをつける
	// validate:"required"は必須パラメータを表す`
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// 引数で与えられた構造体に対してバリデーションを実行する。
	if err := addTask.Validator.Struct(b); err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// t := &entity.Task{
	// 	Title:  b.Title,
	// 	Status: entity.TaskStatusTodo,
	// }

	// err := addTask.Repo.AddTask(ctx, addTask.DB, t)
	t, err := addTask.Service.AddTask(ctx, b.Title)
	if err != nil {
		RespondJson(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	RespondJson(ctx, w, t, http.StatusOK)

}
