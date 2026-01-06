package chat

import (
	"context"
	"web-chat/api/http_model/chat"
)

type Logic interface {
	ResponseStream(ctx context.Context, req *chat.Response) (chat.MessageSteam, error)
	PullModules(ctx context.Context) (*chat.ModelListResp, error)
}
