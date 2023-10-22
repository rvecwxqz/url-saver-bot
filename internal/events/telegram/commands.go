package telegram

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
	"url-saver-bot/internal/clients/telegram"
	"url-saver-bot/internal/storage"
)

const (
	helpCmd      = "/help"
	startCmd     = "/start"
	getCmd       = "/get"
	showAllCmd   = "/show_all"
	removeCmd    = "/remove"
	showTags     = "/show_tags"
	showAllByTag = "/show_all_by_tag"
)

func (p *TgProcessor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%v' from '%v", text, username)

	if isURL(text) {
		return p.addPage(text, username, chatID)
	}

	cmd := strings.Split(text, " ")[0]

	switch cmd {
	case getCmd:
		return p.getPage(username, chatID)
	case helpCmd:
		return p.sendHelp(chatID)
	case startCmd:
		return p.start(chatID)
	case showAllCmd:
		return p.showAll(chatID, username)
	case removeCmd:
		return p.removePage(username, chatID, text)
	case showTags:
		return p.showTags(username, chatID)
	case showAllByTag:
		return p.showAllByTag(username, chatID, text)
	default:
		return p.tgClient.SendMessage(chatID, fmt.Sprintf("%v: %v", unknownCommandMessage, cmd))
	}
}

func (p *TgProcessor) addPage(pageURL string, userName string, chatID int) error {
	sendMsg := NewMessageSender(chatID, p.tgClient)

	page := &storage.Page{
		URL:      pageURL,
		Tags:     "",
		UserName: userName,
		Created:  time.Now(),
	}

	err := p.storage.Save(p.ctx, page)
	var e *storage.AlreadyExistsError
	if errors.As(err, &e) {
		return sendMsg(alreadyExistsMessage)
	} else if err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	p.tagWorker.AppendPage(*page)

	return sendMsg(SavedMessage)
}

func (p *TgProcessor) getPage(userName string, chatID int) error {
	page, err := p.storage.Pick(p.ctx, userName)
	var e *storage.NoResultError
	if errors.As(err, &e) {
		return p.tgClient.SendMessage(chatID, NoSavedPagesMessage)
	} else if err != nil {
		return fmt.Errorf("can't pick URL from storage: %w", err)
	}

	err = p.tgClient.SendMessage(chatID, page.URL)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	if err = p.storage.Remove(p.ctx, page); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

func (p *TgProcessor) removePage(userName string, chatID int, text string) error {
	splitArray := strings.Split(text, " ")
	if len(splitArray) < 2 {
		return p.tgClient.SendMessage(chatID, noLinkMessage)
	}
	URL := splitArray[1]
	if !isURL(URL) {
		return p.tgClient.SendMessage(chatID, noLinkMessage)
	}

	page := storage.Page{
		URL:      URL,
		UserName: userName,
	}
	if err := p.storage.Remove(p.ctx, &page); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return p.tgClient.SendMessage(chatID, pageRemovedMessage)
}

func (p *TgProcessor) showAll(chatID int, userName string) error {
	pages, err := p.storage.PickAll(p.ctx, userName)
	if err != nil {
		return fmt.Errorf("can't get pages: %w", err)
	}

	if len(pages) == 0 {
		return p.tgClient.SendMessage(chatID, NoSavedPagesMessage)
	}
	var out string
	for i, v := range pages {
		out += fmt.Sprintf("\n%v. %v", i+1, v.URL)
	}

	return p.tgClient.SendMessage(chatID, out)
}

func (p *TgProcessor) showTags(userName string, chatID int) error {
	tags, err := p.storage.SelectTags(p.ctx, userName)
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		p.tgClient.SendMessage(chatID, noTagsMessage)
		return nil
	}

	return p.tgClient.SendTags(chatID, tags)
}

func (p *TgProcessor) showAllByTag(userName string, chatID int, text string) error {
	tag := strings.Join(strings.Split(text, " ")[1:], " ")
	urls, err := p.storage.SelectByTag(p.ctx, tag, userName)
	if err != nil {
		return err
	}
	if len(urls) == 0 {
		p.tgClient.SendMessage(chatID, noURLsForTagMessage)
		return nil
	}

	message := fmt.Sprintf("%v:\n%v", tag, strings.Join(urls, "\n"))

	return p.tgClient.SendMessage(chatID, message)

}

func (p *TgProcessor) sendHelp(chatID int) error {
	return p.tgClient.SendMessage(chatID, helpMessage)
}

func (p *TgProcessor) start(chatID int) error {
	return p.tgClient.SendMessage(chatID, helloMessage)
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isURL(text string) bool {
	path, err := url.ParseRequestURI(text)
	if err == nil && strings.ContainsAny(path.Host, ".") {
		return true
	}
	return false
}
