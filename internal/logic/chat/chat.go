package chat

import "web-chat/api/http_model/chat"

type Logic interface {
	ResponseStream(req *chat.Response) (chat.MessageSteam, error)
	PullModules() (*chat.ModelListResp, error)
}
