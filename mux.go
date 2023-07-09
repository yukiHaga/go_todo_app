package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/yukiHaga/go_todo_app/clock"
	"github.com/yukiHaga/go_todo_app/config"
	"github.com/yukiHaga/go_todo_app/handler"
	"github.com/yukiHaga/go_todo_app/service"
	"github.com/yukiHaga/go_todo_app/store"
)

// 戻り値を*http.ServeMux型の値ではなくて、
// http.HandlerはServeHTTPを実装していればOK
// *http.ServeMuxは、ServeHTTPを実装しているから、戻り値でhttp.Handlerでも整合性がある
// 戻り値を*http.ServeMux型ではなくてhttp.Handlerにしておくことで、内部実装に依存しない関数シグネチャになる
// NewMux関数が返すルーティングでは、HTTPサーバーが稼働中かを確認するための/healthエンドポイントを一つ宣言しておく
// コンテナ実行環境の多くでは、コンテナをいつ再起動するかの判断条件として指定されたエンドポイントをポーリングするルールがある
// NewMuxを定義することで、muxライブラリの実装を内部に隠蔽できる。他のファイルのコードはNewMuxと依存することになるので、
// muxと直接依存することがなくなる。その結果、muxライブラリの変更がしやすくなる。muxライブラリが直接他のファイルのコードと依存していたら、
// muxライブラリを変更するのがとてもめんどくさい。今回の場合は変更が一箇所に集中できてる
func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()

	// このNewを使うことで、クライアント側でmysqlを指定するってことは一応なくなるのか
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	r := &store.Repository{Clocker: clock.RealClocker{}}
	// ここで依存性を注入している
	addTaskHandler := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: r},
		Validator: v,
	}
	// ハンドラーのハンドラーファンクションを登録する
	mux.Post("/tasks", addTaskHandler.ServeHTTP)

	listTasksHandler := &handler.ListTasks{
		Service: &service.ListTasks{DB: db, Repo: r},
	}
	mux.Get("/tasks", listTasksHandler.ServeHTTP)

	// ここは依存性を注入している
	registerUserHandler := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: r},
		Validator: v,
	}
	mux.Post("/register", registerUserHandler.ServeHTTP)

	return mux, cleanup, nil
}
