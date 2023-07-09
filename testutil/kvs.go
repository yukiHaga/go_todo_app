package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
)

// テスト環境ごとに接続を変更するRedisのテストヘルパー
func OpenRedisForTest(t *testing.T) *redis.Client {
	t.Helper()

	host := "127.0.0.1"
	port := 36379

	// LookupEnv は、キーで指定された環境変数の値を取得します。
	// その変数が環境に存在する場合、その値 (空でもよい) が返され、ブール値は true になります。
	// そうでない場合は、返される値は空で、ブール値は偽になります。
	// CIの場合は、portを変える
	if _, defined := os.LookupEnv("CI"); defined {
		port = 6379
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})

	// Pingは、RedisクライアントがRedisサーバーに対してPINGコマンドを送信するメソッド。
	// Err()は、前のコマンドのエラーを返すメソッド
	// RedisクライアントがRedisサーバーにPINGコマンドを送信し、その結果をエラーとして取得しています。
	// エラーがnilでない場合、Redisクライアントが正常にRedisサーバーに接続できなかったことを示します。
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("failed to connect redis: %s", err)
	}

	return client
}
