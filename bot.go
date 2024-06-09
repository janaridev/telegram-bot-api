package tgbotapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/janaridev/telegram-bot-api/types"
)

type BotAPI struct {
	Token  string
	client *http.Client
}

func NewBotAPI(token string) *BotAPI {
	return &BotAPI{
		Token:  token,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (b *BotAPI) getAPIURL(method string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", b.Token, method)
}

func (b *BotAPI) GetUpdates(offset, timeout int) ([]types.Update, error) {
	apiURL := b.getAPIURL("getUpdates")
	values := url.Values{}
	values.Set("offset", fmt.Sprintf("%d", offset))
	values.Set("timeout", fmt.Sprintf("%d", timeout))

	resp, err := b.client.Get(apiURL + "?" + values.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Ok     bool           `json:"ok"`
		Result []types.Update `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Result, nil
}

func (b *BotAPI) SendMessage(chatID int, text string) (*types.SendMessageResponse, error) {
	apiURL := b.getAPIURL("sendMessage")
	values := url.Values{}
	values.Set("chat_id", fmt.Sprintf("%d", chatID))
	values.Set("text", text)

	resp, err := b.client.PostForm(apiURL, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result types.SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
