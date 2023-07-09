package auth

// ブランクインポートしないと、このパッケージ使ってないやんってエラーが出る
// 本当はgo:embedで使っている
import (
	"context"
	_ "embed"
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/yukiHaga/go_todo_app/clock"
	"github.com/yukiHaga/go_todo_app/entity"
)

// 秘密鍵と公開鍵を変数に埋め込んどく

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

// ファイルを読み込んだだけだと、鍵ファイルはただのバイト配列でしかない。
// ファイルを読み込んで、鍵として、データを構築する必要がある
// リクエストを処理するたびに同じ内容のバイト配列から「鍵」を生成する必要はない。
// アプリケーション起動時に「鍵」として読み込んだデータを保持するauth.JWTer型を定義する
// auth.JWTer型には、作成したJWTをキーバリューストアに保存するauth.Storeインターフェースのフィールドも定義しておく
type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	Store                 Store
	Clocker               clock.Clocker
}

//go:ggenerate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

// StoreをDIできる構造にしておく
// こうすることで、テストも書きやすくなる。あと、JWTerがStoreに直接依存しない構造にできる
// JWTerもStoreのどちらもインタフェースに準拠している。
// インターフェースを挟んでいるので、JWTerもStoreの変更を気にせずに変更・追加できるし、Store型もJWTerへの影響を気にせずに変更できる
func NewJWTer(s Store) (*JWTer, error) {
	j := &JWTer{Store: s}

	privKey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: private key: %w", err)
	}

	pubkey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: public key: %w", err)
	}

	// JWTのクレームには発行時刻などを入れるので、JWTer.Clockerに入れとく
	j.PrivateKey = privKey
	j.PublicKey = pubkey
	j.Clocker = clock.RealClocker{}
	return j, nil
}

func parse(rawKey []byte) (jwk.Key, error) {
	// ParseKeyを使って、鍵の情報が含まれるバイト列から、jwxパッケージで利用可能なjwk.Keyインターフェースを満たす型を取得する
	key, err := jwk.ParseKey(rawKey, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}

	return key, nil
}
