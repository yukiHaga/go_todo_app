package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// http.Server型をラップした独自のServer型を定義する
type Server struct {
	srv *http.Server
	l   net.Listener
}

// 動的に選択したポートをリッスンするために、net.Listener型の値を引数で受け取る
func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		// マルチプレクサを使用する場合、http.ServerのHandlerフィールドに作成したマルチプレクサを登録する
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

// Server型のRunメソッドを実装する
// この中にHTTPハンドラーの定義は実装しなかった。main.goのrunだとやっていたのに
func (s *Server) Run(ctx context.Context) error {
	// 空のコンテキストに機能を追加した
	// リストされたシグナルのいずれかが到着したとき, 返されたstop関数が呼ばれたとき, または親コンテキストのDoneチャネルが閉じられたときのうち
	// どれかが起こったら、doneとマークされた（そのDoneチャネルが閉じられた）親コンテキストのコピーを返すのか。
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt, os.Kill)
	// defer文でリソースリークを防げる(プログラムが割り当てたリソースを解放していない状態)
	defer stop()

	log.Println("server start")
	// ListenAndServeメソッドではなくて、Serveメソッドに変更する
	// ServeはListener l上で着信コネクションを受け付け、それぞれに対して新しいサービス・ゴルーチンを作成する。
	// 新しいサービスゴルーチンを生成します。サービスゴルーチンはリクエストを読み、srv.Handlerを呼び出して返信します。
	go s.srv.Serve(s.l)

	<-ctx.Done()
	// 返された子ContextのDone チャネルが閉じられるのは、「期限を過ぎた場合」、「CancelFunc が呼び出された場合」、「親Contextの Done チャネルが閉じられた場合」の3パターン。
	// つまり、五秒をすぎた場合、強制的にチャネルが閉じる
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
