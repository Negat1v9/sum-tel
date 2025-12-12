package aihttpclient

// Types for AI HTTP Client
type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	// Temperature float32   `json:"temperature"`
	// Steam       bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseBody struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type RequestBodyV2 struct {
	Message         string   `json:"message"`
	ParentMessageId string   `json:"parent_message_id,omitempty"`
	FileIds         []string `json:"file_ids,omitempty"`
}

type ResponseBodyV2 struct {
	Message string `json:"message"`
	ID      string `json:"id"`
	// FinishReason map[string]any `json:"finish_reason"`
}
