package service

import (
	"context"
	"fmt"

	"github.com/yukiHaga/go_todo_app/store"
)

type Login struct {
	DB             store.Queryer
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

// このpwは平文
func (l *Login) Login(ctx context.Context, name, pw string) (string, error) {
	// userの構造体を取得
	u, err := l.Repo.GetUser(ctx, l.DB, name)
	if err != nil {
		return "", fmt.Errorf("failed to list: %w", err)
	}

	// ユーザー構造体のハッシュ化されたパスワードとリクエストで送られてきた平文のパスワードを比較
	if err := u.ComparePassword(pw); err != nil {
		return "", fmt.Errorf("wrong password: %w", err)
	}

	// パスワードがあっているならこの処理が実行される
	// 署名つきのjwtを生成できる
	jwt, err := l.TokenGenerator.GenerateToken(ctx, *u)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	// バイトスライスを文字列になおした
	return string(jwt), nil
}
