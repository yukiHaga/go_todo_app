package main

import "net/http"

// 戻り値を*http.ServeMux型の値ではなくて、
// http.HandlerはServeHTTPを実装していればOK
// *http.ServeMuxは、ServeHTTPを実装しているから、戻り値でhttp.Handlerでも整合性がある
// 戻り値を*http.ServeMux型ではなくてhttp.Handlerにしておくことで、内部実装に依存しない関数シグネチャになる
// NewMux関数が返すルーティングでは、HTTPサーバーが稼働中かを確認するための/healthエンドポイントを一つ宣言しておく
// コンテナ実行環境の多くでは、コンテナをいつ再起動するかの判断条件として指定されたエンドポイントをポーリングするルールがある
func NewMux() http.Handler {
	// マルチプレクサを作成
	mux := http.NewServeMux()
	// mux.HandleFuncメソッドで、マルチプレクサにURLと対応する処理(ハンドラ)を登録する
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
