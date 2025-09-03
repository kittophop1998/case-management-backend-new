package utils

import (
	"context"
	"net/http"
)

func SetHeadersFormContext(ctx context.Context, req *http.Request, keys []CtxKey) {
	for _, key := range keys {
		if val, ok := ctx.Value(key).(string); ok && val != "" {
			req.Header.Set(string(key), val)
		} else {

		}
	}
}
