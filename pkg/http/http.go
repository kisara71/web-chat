package http

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RequestHandler struct {
	commonClient *http.Client
	sseClient    *http.Client
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		commonClient: &http.Client{Timeout: time.Second * 5},
		sseClient:    &http.Client{Timeout: 0},
	}
}

func (r *RequestHandler) DoCommon(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := r.commonClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return nil, fmt.Errorf("http %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}
	return resp, nil
}

type SSEEvent struct {
	ID    string
	Data  string
	Event string
}
type SSEReader struct {
	resp   *http.Response
	reader *bufio.Reader
}

func (sr *SSEReader) Close() error {
	return sr.resp.Body.Close()
}
func (sr *SSEReader) Next() (SSEEvent, bool, error) {
	var ev SSEEvent
	var data strings.Builder

	flush := func() (SSEEvent, bool, error) {
		if data.Len() > 0 {
			ev.Data = data.String()
			if strings.TrimSpace(ev.Data) == "[DONE]" {
				return SSEEvent{}, false, nil
			}
			return ev, true, nil
		}
		return SSEEvent{}, true, nil
	}

	for {
		line, err := sr.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// 尝试把最后一个事件 flush 掉
				if data.Len() > 0 || ev.Event != "" || ev.ID != "" {
					out, ok, ferr := flush()
					if ferr != nil {
						return SSEEvent{}, false, ferr
					}
					if ok {
						return out, true, nil
					}
				}
				return SSEEvent{}, false, nil
			}
			return SSEEvent{}, false, err
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			out, ok, ferr := flush()
			if ferr != nil {
				return SSEEvent{}, false, ferr
			}
			if !ok {
				return SSEEvent{}, false, nil
			}
			// out 为空表示这个事件没 data，继续读下一个
			if out.Data == "" && out.Event == "" && out.ID == "" {
				ev = SSEEvent{}
				data.Reset()
				continue
			}
			return out, true, nil
		}

		// 注释/心跳
		if strings.HasPrefix(line, ":") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "event:"):
			ev.Event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		case strings.HasPrefix(line, "id:"):
			ev.ID = strings.TrimSpace(strings.TrimPrefix(line, "id:"))
		case strings.HasPrefix(line, "data:"):
			v := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if data.Len() > 0 {
				data.WriteByte('\n')
			}
			data.WriteString(v)
		default:
		}
	}
}
func (r *RequestHandler) DoSSE(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*SSEReader, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := r.sseClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return nil, fmt.Errorf("http %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	return &SSEReader{
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}, nil
}
