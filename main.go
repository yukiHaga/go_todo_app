package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/yukiHaga/go_todo_app/config"
)

// run関数のctxはリクエストのctxではない
// run関数に書いている処理を直接mainに書くと、テストがしづらく、以下の問題がある
// - テスト完了後に終了する術がない(サーバーだから常に起動している)
// - 出力を検証しにくい(main関数だから戻り値がない)
// - 異常時にos.Exit関数が呼ばれて、直ちに終了してしまう。
// - ポート番号が固定されているので、サーバーを起動したままテストを実行しようとすると、ポートが利用できずにテストが失敗する。
func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.AppPort))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.AppPort, err)
	}

	mux := NewMux()
	s := NewServer(l, mux)
	return s.Run(ctx)
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Println("error: server end")
		os.Exit(1)
	}
}
