package logat

import (
	"context"
)

type contextInfo struct {
	UserID uint64 `json:"user_id"`
}

func getContext(ctx context.Context) contextInfo {
	userID, _ := ctx.Value(contextKeyUserID).(uint64)

	contextInfoObj := contextInfo{
		UserID: userID,
	}
	return contextInfoObj
}
