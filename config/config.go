// 環境変数を読み込んで構造体にマッピングするために、configパッケージを定義した
package config

// envパッケージは環境変数を構造体をマッピングするためのもの
import "github.com/caarlos0/env/v9"

// 各環境変数をフィールドとして持つ構造体
type Config struct {
	Env        string `env:"TODO_ENV" envDefault:"dev"`
	AppPort    int    `env:"APP_PORT" envDefault:"80"`
	DBHost     string `env:"TODO_DB_HOST" envDefault:"127.0.0.1"`
	DBPort     int    `env:"TODO_DB_PORT" envDefault:"33306"`
	DBUser     string `env:"TODO_DB_USER" envDefault:"todo"`
	DBPassword string `env:"TODO_DB_PASSWORD" envDefault:"todo"`
	DBName     string `env:"TODO_DB_NAME" envDefault:"todo"`
	RedisHost  string `env:"TODO_REDIS_HOST" envDefault:"127.0.0.1"`
	RedisPort  int    `env:"TODO_REDIS_PORT" envDefault:"36379"`
}

func New() (*Config, error) {
	cfg := &Config{}
	// Parse は `env` タグを含む構造体を解析し、その値を環境変数からロードする。
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
