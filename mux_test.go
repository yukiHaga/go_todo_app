package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// ルーティングが意図通りかテストしている
// TestNewMux関数ではhttptestパッケージを使って、ServeHTTP関数の引数に渡すためのモックを作成している
func TestNewMux(t *testing.T) {
	// ResponseWriterインターフェースを満たす*ResponseRecorder型の値を取得します
	w := httptest.NewRecorder()
	// 本当にリクエストだけ作成している。リクエストを実行したわけではない
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	sut := NewMux()
	// *ResponseRecorder型の値をServeHTTP関数に渡した後にResultメソッドを実行すると、
	// クライアントが受け取るレスポンス内容が含まれるhttp.Response型の値を取得できる
	// この処理はハンドラーに定義してある関数を自分で呼び出しているだけか
	sut.ServeHTTP(w, r)
	resp := w.Result()
	// Cleanup は、テスト（またはサブテスト）とそのすべてのサブテストが完了したときに呼び出される関数を登録します
	// Closeすることで、ボディが占有しているメモリのリソースを解放している
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusOK {
		t.Errorf("want status code")
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
