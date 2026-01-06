package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	httpmodel "web-chat/api/http_model/chat"
	"web-chat/internal/logic/chat"
	"web-chat/internal/svc"
	http2 "web-chat/pkg/http"
	"web-chat/pkg/logger"
	"web-chat/pkg/utils"
)

type logicImpl struct {
	svcCtx  *svc.Context
	utils   *utils.Utils
	urls    *urls
	authKey string
	headers map[string]string
}

func NewChatLogic(svcCtx *svc.Context) (chat.Logic, error) {
	key := os.Getenv("OPENAI_KEY")
	if key == "" {
		return nil, errors.New("empty openai key")
	}
	baseURL := svcCtx.Config.LLMRequestConf.OpenAI.BaseURL
	if baseURL == "" {
		return nil, errors.New("empty openai base url")
	}
	headers := map[string]string{
		"Authorization": "Bearer " + key,
	}
	return &logicImpl{
		svcCtx:  svcCtx,
		utils:   svcCtx.Utils,
		urls:    newURLs(baseURL),
		authKey: "Bearer " + key,
		headers: headers,
	}, nil
}
func (l *logicImpl) ResponseStream(ctx context.Context, req *httpmodel.Response) (httpmodel.MessageSteam, error) {
	var (
		sr  *http2.SSEReader
		err error
	)
	if len(req.Messages) == 0 {
		return nil, errors.New("empty message")
	}
	req.Stream = true
	bodyByte, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	sr, err = l.utils.RequestHandler.DoSSE(ctx, http.MethodPost, l.urls.Completion, bytes.NewReader(bodyByte), l.headers)
	if err != nil {
		return nil, err
	}
	return newOpenAIResponsesStream(sr), nil
}

func (l *logicImpl) PullModules(ctx context.Context) (*httpmodel.ModelListResp, error) {
	var (
		res = new(httpmodel.ModelListResp)
		err error
	)

	resp, err := l.utils.RequestHandler.DoCommon(
		ctx,
		http.MethodGet,
		l.urls.ModelList,
		nil,
		l.headers,
	)

	if err != nil {
		logger.L().Errorf("pull modules error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
