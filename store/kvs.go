package store

import (
	"context"
	"fmt"
	"time"

	// go-redisはredisクライアントとしてよく使われている
	"github.com/redis/go-redis/v9"
	"github.com/yukiHaga/go_todo_app/config"
	"github.com/yukiHaga/go_todo_app/entity"
)

// Redisクライアントを初期化するコンストラクタ
func NewKVS(ctx context.Context, cfg *config.Config) (*KVS, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})

	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &KVS{Cli: cli}, nil
}

// RedisClientのラッパー
// おそらくkey-value-dbにredisを採用しなくなっても良いように、抽象度を上げている
// ただ、もしそうだとしたら、環境変数で台無し感はあるけど
type KVS struct {
	Cli *redis.Client
}

// アクセストークンのID(JWTのClaimのjti属性)をキーとして、値にはユーザーのIDを保存する設計にしている
// userIDをint64ではなくて、entity.UserIDにすることで、データの誤取り扱いを防ぐ
func (k *KVS) Save(ctx context.Context, key string, userID entity.UserID) error {
	// 型を変えているのか。
	id := int64(userID)
	// 30分がキーの有効ってことか
	// 有効期限が過ぎると、Redisは自動的にキーを削除する
	// Redisにはkeyとid(value)を保存している
	// キーが削除されると、関連するバリューも自動的に削除される
	return k.Cli.Set(ctx, key, id, 30*time.Minute).Err()
}

// Redisからデータをロードしている
func (k *KVS) Load(ctx context.Context, key string) (entity.UserID, error) {
	id, err := k.Cli.Get(ctx, key).Int64()
	if err != nil {
		return 0, fmt.Errorf("failed to get by %q: %w", key, ErrNotFound)
	}
	return entity.UserID(id), nil
}
