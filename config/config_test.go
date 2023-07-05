package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	wantPort := 3333

	// t.Setenvはos.Setenv(key, value)を呼び出し、テスト後にCleanupを使用して環境変数を元の値に戻します。
	// os.Setenvは、キーで指定された環境変数の値を設定する。エラーがあればエラーを返す。
	t.Setenv("APP_PORT", fmt.Sprint(wantPort))

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config %v", err)
	}

	if got.AppPort != wantPort {
		t.Errorf("want %d but %d", wantPort, got.AppPort)
	}

	wantEnv := "dev"

	if got.Env != wantEnv {
		t.Errorf("want %s, but %s", wantEnv, got.Env)
	}
}
