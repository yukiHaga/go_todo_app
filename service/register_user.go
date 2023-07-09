package service

import (
	"context"
	"fmt"

	"github.com/yukiHaga/go_todo_app/entity"
	"github.com/yukiHaga/go_todo_app/store"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	DB   store.Execer
	Repo UserRegister
}

func (r *RegisterUser) RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error) {
	// 与えられたコストでパスワードの bcrypt ハッシュを返します。指定されたコストが MinCost より小さい場合は、代わりに DefaultCost が設定されます。
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	u := &entity.User{
		Name:     name,
		Password: string(pw),
		Role:     role,
	}

	// ちゃんとハッシュ化されてい文字列が表示されていた。
	// つまりstringかける前は、ハッシュ化した文字列のバイナリってことか
	// log.Println("=========")
	// log.Println(string(pw))
	// log.Println("=========")

	if err := r.Repo.RegisterUser(ctx, r.DB, u); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return u, nil
}
