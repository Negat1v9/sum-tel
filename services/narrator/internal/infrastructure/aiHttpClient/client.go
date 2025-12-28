package aihttpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Negat1v9/sum-tel/services/narrator/internal/domain"
	"github.com/Negat1v9/sum-tel/shared/config"
)

type Client struct {
	c       http.Client
	baseUrl string

	systemPromt string
	token       string
	model       string
	isStreaming bool
}

func NewClient(cfg *config.AiClientCfg, isStreaming bool) *Client {
	return &Client{
		c:           http.Client{},
		baseUrl:     cfg.BaseUrl,
		systemPromt: cfg.SystemPrompt,
		token:       cfg.Token,
		model:       cfg.Model,
		isStreaming: isStreaming,
	}
}

// sends a request to the AI service to aggregate raw messages into a summarized response
func (c *Client) DoAggregation(ctx context.Context, msgs []domain.RawMessage) (*domain.AggregationResponse, int, error) {
	mn := "Client.DoAggregation"
	type Req struct {
		Messages []domain.RawMessage `json:"messages"`
	}
	bMsgs, err := json.Marshal(&Req{Messages: msgs})
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", mn, err)
	}
	body := RequestBody{
		Model: c.model,
		Messages: []Message{
			{
				Role:    "user",
				Content: string(bMsgs),
			},
			{
				Role:    "system",
				Content: c.systemPromt,
			},
		},
	}

	bBody, err := json.Marshal(&body)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", mn, err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseUrl, bytes.NewReader(bBody))
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", mn, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", mn, err)
	}

	defer resp.Body.Close()

	var respBody ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, 0, fmt.Errorf("%s unmarshal ResponseBody: %w", mn, err)
	}

	if len(respBody.Choices) == 0 {
		return nil, respBody.Usage.TotalTokens, fmt.Errorf("%s no choices", mn)
	}

	var aggregation domain.AggregationResponse
	err = json.Unmarshal([]byte(clearJson(respBody.Choices[0].Message.Content)), &aggregation)
	if err != nil {
		return nil, respBody.Usage.TotalTokens, fmt.Errorf("%s unmarshal AggregationResponse: %w", mn, err)
	}

	return &aggregation, respBody.Usage.TotalTokens, nil
}

func clearJson(s string) string {
	s, _ = strings.CutPrefix(s, "```json")
	s, _ = strings.CutSuffix(s, "```")
	return s
}
