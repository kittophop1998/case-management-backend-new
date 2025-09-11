package utils

import (
	"context"
	"net/http"
)

func SetHeadersFormContext(ctx context.Context, req *http.Request, keys []CtxKey) {
	for _, key := range keys {
		if val := ctx.Value(key); val != nil {
			if str, ok := val.(string); ok && str != "" {
				req.Header.Set(string(key), str)
			}
		}
	}
}
