package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/yukiHaga/go_todo_app/entity"
	"github.com/yukiHaga/go_todo_app/store"
	"github.com/yukiHaga/go_todo_app/testutil"
)

func TestAddTask(t *testing.T) {
	t.Parallel()
	type want struct {
		status  int
		rspFile string
	}

	cases := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_req_rsp.json.golden",
			},
		},
	}

	for n, c := range cases {
		c := c
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			// httpパッケージを使うことで、HTTPリクエストやHTTPサーバーをシミュレートしてテストを行うことができる
			// 以下の場合はHTTPリクエストをシミュレートしている
			// ResponseWriterインターフェースを満たす*ResponseRecorder型の値を取得します
			// レコーダーでレスポンスを記録している
			w := httptest.NewRecorder()

			// 本当にリクエストだけ作成している。リクエストを実行したわけではない
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, c.reqFile)),
			)

			// *ResponseRecorder型の値をServeHTTP関数に渡した後にResultメソッドを実行すると、
			// クライアントが受け取るレスポンス内容が含まれるhttp.Response型の値を取得できる
			// この処理はハンドラーに定義してある関数を自分で呼び出しているだけか
			sut := AddTask{
				Store: &store.TaskStore{
					Tasks: map[entity.TaskID]*entity.Task{},
				},
				Validator: validator.New(),
			}
			sut := NewAddTask()
			sut.ServeHTTP(w, r)
			resp := w.Result()
			testutil.AssertResponse(t, resp, c.want.status, testutil.LoadFile(t, c.want.rspFile))
		})
	}

}
