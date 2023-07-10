package auth

// ブランクインポートしないと、このパッケージ使ってないやんってエラーが出る
// 本当はgo:embedで使っている
import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
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

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

// StoreをDIできる構造にしておく
// こうすることで、テストも書きやすくなる。あと、JWTerがStoreに直接依存しない構造にできる
// JWTerもStoreのどちらもインタフェースに準拠している。
// インターフェースを挟んでいるので、JWTerもStoreの変更を気にせずに変更・追加できるし、Store型もJWTerへの影響を気にせずに変更できる
func NewJWTer(s Store, c clock.Clocker) (*JWTer, error) {
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
	j.Clocker = c
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

const (
	RoleKey     = "role"
	UserNameKey = "user_name"
)

func (j *JWTer) GenerateToken(ctx context.Context, u entity.User) ([]byte, error) {
	// JWTの中身をビルダーパターンで組み立てる
	// Build() が呼び出されたときに Builder がトークンを構築します。
	token, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/yukiHaga/go_todo_app`).
		Subject("access_token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(30*time.Minute)).
		Claim(RoleKey, u.Role).
		Claim(UserNameKey, u.Name).
		Build()
	if err != nil {
		return nil, fmt.Errorf("GetToken: failed to builed token: %w", err)
	}
	// ユーザーidをRedisに登録
	if err := j.Store.Save(ctx, token.JwtID(), u.ID); err != nil {
		return nil, err
	}

	// Signは、コンパクトな形式でシリアライズされた署名付きJWTトークンを作成する便利な関数です。
	// トークンと、生の鍵またはjwk.Keyと、トークンの署名に必要なアルゴリズムを受け取る
	// WithKeyは多目的オプションです。jwt.Sign、jwt.Parse（およびその兄弟）、jwt.Serializerメソッドのいずれにも使用できます。
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, err
	}
	return signed, nil
}

func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	// ParseRequestを使うことで、httpリクエストから、JWTであるjwt.Tokenインターフェースを満たす型の値を取得できる
	// jwt.WithKey関数は署名を検証するアルゴリズムと利用する鍵を指定している
	// jwt.WithValidate関数を使うことで、検証は無視している。これはDIをしている*auth.JWTer.Clockerフィールドをベースに検証を行うためである
	// それだと、改ざんしている場合に見抜けないけど大丈夫か？って感じがするな。機嫌だけで検証するのは危険すぎる
	token, err := jwt.ParseRequest(
		r,
		jwt.WithKey(jwa.RS256, j.PublicKey),
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, err
	}

	// Validateはクレームが正しいかどうかを検証する
	// 時刻の検証で使っているそう
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, fmt.Errorf("GetToken: failed to validate token: %w", err)
	}

	// JWTが共有メモリ上に存在しているかを一応確認している
	// Redisから削除して手動でexpireさせていることもありうる
	// JWTの jti (JWT ID) は、JWTの一意の識別子を表すクレーム（Claim）の1つである。
	// JwtIDで取得できる
	if _, err := j.Store.Load(ctx, token.JwtID()); err != nil {
		return nil, fmt.Errorf("GetToken: %q expired* %w", token.JwtID(), err)
	}
	return token, nil
}
