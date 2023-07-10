package auth

import (
	"bytes"
	"context"
	"testing"

	"github.com/yukiHaga/go_todo_app/clock"
	"github.com/yukiHaga/go_todo_app/entity"
	"github.com/yukiHaga/go_todo_app/testutil/fixture"
)

func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN PUBLIC KEY-----")
	// Containsはサブスライスが b 内にあるかどうかを報告します。
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("want %s, but got %s", want, rawPubKey)
	}

	want = []byte("-----BEGIN PRIVATE KEY-----")
	if !bytes.Contains(rawPrivKey, want) {
		t.Errorf("want %s, but got %s", want, rawPrivKey)
	}

}

func TestJWTer_GenerateToken(t *testing.T) {
	ctx := context.Background()
	moq := &StoreMock{}
	wantId := entity.UserID(20)
	u := fixture.User(&entity.User{ID: wantId})

	// Saveが呼ばれた時に、登録した関数を実行する。
	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantId {
			t.Errorf("want %d, but got %d", wantId, userID)
		}
		return nil
	}

	sut, err := NewJWTer(moq, clock.RealClocker{})
	if err != nil {
		t.Fatal(err)
	}

	got, err := sut.GenerateToken(ctx, *u)
	if err != nil {
		t.Fatalf("not want err: %v", err)
	}

	if len(got) == 0 {
		t.Errorf("token is empty")
	}
}

// func TestJWTer_GetToken(t *testing.T) {
// 	t.Parallel()

// 	c := clock.FixedClocker{}
// 	want, err := jwt.NewBuilder().
// 		JwtID(uuid.New().String()).
// 		Issuer(`github.com/yukiHaga/go_todo_app`).
// 		Subject("access_token").
// 		IssuedAt(c.Now()).
// 		Expiration(c.Now().Add(30*time.Minute)).
// 		Claim(RoleKey, "test").
// 		Claim(UserNameKey, "test_user").
// 		Build()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// ParseKeyを使って、鍵の情報が含まれるバイト列から、jwxパッケージで利用可能なjwk.Keyインターフェースを満たす型を取得する
// 	pkey, err := jwk.ParseKey(rawPrivKey, jwk.WithPEM(true))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	signed, err := jwt.Sign(want, jwt.WithKey(jwa.RS256, pkey))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	userID := entity.UserID(20)

// 	ctx := context.Background()
// 	moq := &StoreMock{}
// 	moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
// 		return userID, nil
// 	}
// 	sut, err := NewJWTer(moq, c)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	req := httptest.NewRequest(
// 		http.MethodGet,
// 		`https://github.com/yukiHaga`,
// 		nil,
// 	)
// 	req.Header.Set(`Authorization`, fmt.Sprintf(`Bearer %s`, signed))
// 	got, err := sut.GetToken(ctx, req)
// 	if err != nil {
// 		t.Fatalf("want no error, but got %v", err)
// 	}
// 	if !reflect.DeepEqual(got, wasnt) {
// 		t.Errorf("GetToken() got = ")
// 	}
// }
