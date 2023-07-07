package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func RespondJson(ctx context.Context, w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json: charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rsp := ErrResponse{
			// StatusText は、HTTP ステータス・コードのテキストを返します。
			Message: http.StatusText(http.StatusInternalServerError),
		}
		// json.NewEncoder(w) によって、指定された w を出力先とする新しい JSON エンコーダーが作成されます。
		// Encode(rsp) が呼び出され、エンコーダーが rsp を JSON 形式に変換します。
		// 変換された JSON データが w に書き込まれます。
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			log.Printf("write error response error: %v", err)
		}
	}

	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}
