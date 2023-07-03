package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
)

// なんのテストをやっているかわかりづらいね
// コメントを書かないといけないのがだるい
func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	// リスナーの生成に失敗しているかを確認
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error)
	go func() {
		err := run(ctx, l)
		errCh <- err
	}()

	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), "message")
	// どんなポート番号でリッスンしているか確認
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)

	// リクエストが正しく飛ぶか
	if err != nil {
		// テスト失敗のログを出力して処理を継続(フォーマットあり)
		t.Errorf("failed to get %+v", err)
	}
	defer rsp.Body.Close()

	// ボディを参照できるか
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		// テスト失敗のログを出力して処理終了(フォーマットあり)
		t.Fatalf("failed to read body: %v", err)
	}

	// レスポンスボディは想定したものか
	want := fmt.Sprintf("Hello, %s", "message")
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	cancel()
	<-errCh

	// エラーは出ていないか
	if err != nil {
		t.Fatal(err)
	}
}
