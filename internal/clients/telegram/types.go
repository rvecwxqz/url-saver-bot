package telegram

type UpdateRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID            int              `json:"update_id"`
	Message       *IncomingMessage `json:"message,omitempty"`
	CallbackQuery *CallbackQuery   `json:"callback_query,omitempty"`
}

type MessageRequest struct {
	ChatID             int                   `json:"chat_id"`
	Text               string                `json:"text"`
	DisablePagePreview bool                  `json:"disable_web_page_preview"`
	ReplyMarkup        *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type MessageResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
	Result      string `json:"result"`
	ErrorCode   int    `json:"error_code"`
}

type IncomingMessage struct {
	Chat Chat   `json:"chat"`
	From User   `json:"from"`
	Text string `json:"text"`
}

type Chat struct {
	ID int `json:"id"`
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type CallbackQuery struct {
	From    User            `json:"from"`
	Message IncomingMessage `json:"message"`
	Data    string          `json:"data"`
}
