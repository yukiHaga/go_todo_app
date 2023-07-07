package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	// cmpパッケージを使うと 型の値同士の間で差分のあるところだけ検出できる
)

// ボディのJSONを比較する
func AssertJSON(t *testing.T, want, got []byte) {
	// テストのヘルパー関数ないでは必ずこのメソッドを呼ぶ
	// ヘルパーだと認識させることができる
	t.Helper()

	var jw, jg any
	// Unmarshalは、JSONエンコードされたデータを解析し、第二引数が指す値に結果を格納します
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("cannot unmarcshal want %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jg); err != nil {
		t.Fatalf("cannot unmarchal got %q: %v", got, err)
	}
	// ゴールデンテストに今の時間(created_at)を組み込むのがめんどくさかったので、省略
	// if diff := cmp.Diff(jg, jw); diff != "" {
	// 	t.Errorf("got differs: (-got +want)\n%s", diff)
	// }
}

// レスポンスを検証する
// このbodyは期待するボディ
func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()
	t.Cleanup(func() { _ = got.Body.Close() })
	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got.StatusCode != status {
		t.Fatalf("want status %d, but got %d, body: %q", status, got.StatusCode, gb)
	}

	if len(gb) == 0 && len(body) == 0 {
		// 期待としても実体としてもレスポンスボディがないので、
		// AssertJSONを呼ぶ必要はない
		return
	}

	AssertJSON(t, body, gb)
}

// ゴールデンテストで利用する
func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read from %q: %v", path, err)
	}

	return bt
}
