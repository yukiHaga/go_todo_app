package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// run関数のctxはリクエストのctxではない
func run(ctx context.Context, l net.Listener) error {
	// 空のコンテキストに機能を追加した
	// リストされたシグナルのいずれかが到着したとき, 返されたstop関数が呼ばれたとき, または親コンテキストのDoneチャネルが閉じられたときのうち
	// どれかが起こったら、doneとマークされた（そのDoneチャネルが閉じられた）親コンテキストのコピーを返すのか。
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt, os.Kill)
	// defer文でリソースリークを防げる(プログラムが割り当てたリソースを解放していない状態)
	defer stop()

	s := &http.Server{
		// 引数で受け取ったnet.Listenerを利用するので、
		// Addrフィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(5 * time.Second)
			log.Println("sever response")
			fmt.Fprintf(w, "Hello, %s", r.URL.Path[1:])
		}),
	}

	log.Println("server start")
	// ListenAndServeメソッドではなくて、Serveメソッドに変更する
	go s.Serve(l)

	<-ctx.Done()
	// 返された子ContextのDone チャネルが閉じられるのは、「期限を過ぎた場合」、「CancelFunc が呼び出された場合」、「親Contextの Done チャネルが閉じられた場合」の3パターン。
	// つまり、五秒をすぎた場合、強制的にチャネルが閉じる
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// コマンドラインで与えられたパラメータはosパッケージのArgsで受け取ることができる
	// os.Argsはstring型のスライス
	// 0番目の要素には実行したコマンド名が格納される
	// 1番目以降の各要素にコマンドに渡された各引数が格納される
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}

	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)

	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}

	if err := run(context.Background(), l); err != nil {
		log.Println("error: server end")
	}
}
