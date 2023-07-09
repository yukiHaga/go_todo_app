package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	// インポートはしないけどパッケージ内に定義された初期化処理だけをしたいからブランクインポートしている
	// 普通にインポートしちゃうとこのパッケージ使ってないでしょって怒られるから、明示的にブランクインポートしている
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func OpenDBForTest(t *testing.T) *sqlx.DB {
	t.Helper()

	port := 33306
	// CIがあるなら、ポートを切り替える
	// ローカル環境やGitHub Actions上の環境ではポート番号のみが異なるので、ポート番号のみを切り替えています
	if _, defined := os.LookupEnv("CI"); defined {
		port = 3306
	}

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("todo:todo@tcp(127.0.0.1:%d)/todo?parseTime=true", port),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(
		func() { db.Close() },
	)

	return sqlx.NewDb(db, "mysql")
}
