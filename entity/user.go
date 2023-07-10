package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserID int64

type User struct {
	ID        UserID    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password" db:"password"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// パスワード検証ロジック
// ハッシュ化されたパスワードとそれに相当する可能性がある平文とを比較する
// 成功した場合はnilを返す
func (u *User) ComparePassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
}
