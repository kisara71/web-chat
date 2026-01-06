package chat

type ModelListResp struct {
	Data []struct {
		ID        string `json:"id"`
		CreatedAt int64  `json:"created"`
	} `json:"data"`
}

type Response struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream bool `json:"stream"`
}
