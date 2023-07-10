package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/yukiHaga/go_todo_app/entity"
)

func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	u.CreatedAt = r.Clocker.Now()
	u.UpdatedAt = r.Clocker.Now()
	sql := `INSERT INTO users (name, password, role, created_at, updated_at)
	        VALUES (?, ?, ?, ?, ?)`

	result, err := db.ExecContext(ctx, sql, u.Name, u.Password, u.Role, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry)
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = entity.UserID(id)
	return nil
}

func (r *Repository) GetUser(ctx context.Context, db Queryer, name string) (*entity.User, error) {
	u := &entity.User{}
	sql := `SELECT id, name, password, role, created_at, updated_at FROM users WHERE name = ?`
	// GetContextを使うと、クエリの実行結果が設定された構造体を簡単に取得できる
	if err := db.GetContext(ctx, u, sql, name); err != nil {
		return nil, err
	}

	return u, nil
}
