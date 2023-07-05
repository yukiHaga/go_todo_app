// 環境変数を読み込んで構造体にマッピングするために、configパッケージを定義した
package config

import "github.com/caarlos0/env/v9"

type Config struct {
	Env     string `env:"TODO_ENV" envDefault:"dev"`
	AppPort int    `env:"APP_PORT" envDefault:"80"`
}

func New() (*Config, error) {
	cfg := &Config{}
	// Parse は `env` タグを含む構造体を解析し、その値を環境変数からロードする。
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
