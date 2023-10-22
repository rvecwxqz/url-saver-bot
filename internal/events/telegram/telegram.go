package telegram

import (
	"context"
	"fmt"
	"url-saver-bot/internal/clients/telegram"
	"url-saver-bot/internal/events"
	"url-saver-bot/internal/ml/parser"
	"url-saver-bot/internal/storage"
)

type TgProcessor struct {
	tgClient  *telegram.Client
	offset    int
	storage   storage.Storage
	tagWorker *parser.TagWorker
	ctx       context.Context
}

type Meta struct {
	ChatID       int
	UserID       int
	UserName     string
	CallbackData string
}

func New(ctx context.Context, c *telegram.Client, s storage.Storage, maxBufferSize int) *TgProcessor {
	return &TgProcessor{
		tgClient:  c,
		storage:   s,
		tagWorker: parser.NewTagWorker(ctx, s, maxBufferSize),
		ctx:       ctx,
	}
}

func (p *TgProcessor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tgClient.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *TgProcessor) Process(e events.Event) error {
	switch e.Type {
	case events.Message:
		return p.processMessage(e)
	case events.Callback:
		return p.processCallback(e)
	default:
		return NewUnknownTypeError()
	}
}

func (p *TgProcessor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message %w", err)
	}

	if err = p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	return nil
}

func (p *TgProcessor) processCallback(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process callback %w", err)
	}
	cmd := fmt.Sprintf("%v %v", showAllByTag, meta.CallbackData)
	if err = p.doCmd(cmd, meta.ChatID, meta.UserName); err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}
	return nil
}

func meta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, NewUnknownMetaTypeError()
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	switch updType {
	case events.Message:
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserID:   upd.Message.From.ID,
			UserName: upd.Message.From.UserName,
		}
	case events.Callback:
		res.Meta = Meta{
			ChatID:       upd.CallbackQuery.Message.Chat.ID,
			UserID:       upd.CallbackQuery.From.ID,
			UserName:     upd.CallbackQuery.From.UserName,
			CallbackData: upd.CallbackQuery.Data,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	} else if upd.CallbackQuery != nil {
		return events.Callback
	}
	return events.Unknown
}
