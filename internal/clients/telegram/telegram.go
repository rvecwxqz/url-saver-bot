package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	tgHost            = "api.telegram.org"
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
	showTagMessage    = "Here is all your tags:"
)

type Client struct {
	host     string
	basePath string
	client   *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		host:     tgHost,
		basePath: newBasePath(token),
		client:   &http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	u := UpdateRequest{
		Offset: offset,
		Limit:  limit,
	}

	body, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("update request marshalling error: %w", err)
	}

	data, err := c.doRequest(getUpdatesMethod, body)
	if err != nil {
		return nil, err
	}

	var resp UpdateResponse

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	m := MessageRequest{
		ChatID:             chatID,
		Text:               text,
		DisablePagePreview: true,
	}

	body, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("message marshalling error: %w", err)
	}

	_, err = c.doRequest(sendMessageMethod, body)
	if err != nil {
		return fmt.Errorf("send message error: %w", err)
	}

	return nil
}

func (c *Client) SendTags(chatID int, tags []string) error {
	messageRequest := MessageRequest{
		ChatID:             chatID,
		Text:               showTagMessage,
		DisablePagePreview: false,
		ReplyMarkup:        createReplyMarkup(tags),
	}

	body, err := json.Marshal(messageRequest)
	if err != nil {
		return fmt.Errorf("can't marshal json: %w", err)
	}

	body, err = c.doRequest(sendMessageMethod, body)
	if err != nil {
		return fmt.Errorf("send message error: %w", err)
	}
	var message MessageResponse
	err = json.Unmarshal(body, &message)

	return nil
}

func createReplyMarkup(tags []string) *InlineKeyboardMarkup {
	countInRow := 5
	countRows := int(math.Ceil(float64(len(tags)) / float64(countInRow)))
	residual := len(tags)

	nextTag := tag(tags)
	buttons := make([][]InlineKeyboardButton, 0, countRows)
	for i := 0; i < countRows; i++ {
		countInRow = int(math.Min(float64(countInRow), float64(residual)))
		buttonArray := make([]InlineKeyboardButton, 0, countInRow)
		for j := 0; j < countInRow; j++ {
			tag := strings.TrimSpace(nextTag())
			button := InlineKeyboardButton{
				Text:         tag,
				CallbackData: tag,
			}
			buttonArray = append(buttonArray, button)
		}
		residual -= countInRow
		buttons = append(buttons, buttonArray)
	}

	return &InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func tag(tags []string) func() string {
	i := -1
	return func() string {
		i++
		return tags[i]
	}
}
func (c *Client) doRequest(method string, body []byte) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, NewRequestError(err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, NewRequestError(err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, NewRequestError(err)
	}

	return body, nil

}
