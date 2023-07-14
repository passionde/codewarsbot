package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (b *Bot) handlerStart(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hey! Send me a user's nickname in CodeWars and I'll collect his account statistics")
	b.bot.Send(msg)
}

func (b *Bot) handlerStatistics(update *tgbotapi.Update) {
	username := update.Message.Text
	chatId := update.Message.Chat.ID

	info, ok := b.cache.Get(username)
	if ok {
		msg := tgbotapi.NewMessage(chatId, info.data["all"])
		msg.ReplyMarkup = b.generatePaginate(username, info.keys)
		b.bot.Send(msg)
		return
	}

	response, err := b.store.AllInfo(username)
	if err != nil {
		switch err.Error() {
		case "not found":
			b.bot.Send(tgbotapi.NewMessage(chatId, fmt.Sprintf("Пользователь \"%s\" не найден", username)))
		default:
			b.bot.Send(tgbotapi.NewMessage(chatId, "Ooops, что-то пошло не так) Попробуйте позже"))
			b.logger.Error(err.Error(), username)
		}
		return
	}

	data, keys := generateMessage(response)
	b.cache.Set(username, data, keys)

	msg := tgbotapi.NewMessage(chatId, data["all"])
	msg.ReplyMarkup = b.generatePaginate(username, keys)
	b.bot.Send(msg)
	return
}

func (b *Bot) handlerSetPage(update *tgbotapi.Update) {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
	b.bot.Request(callback)

	parseData := strings.Split(update.CallbackQuery.Data, ":")
	if len(parseData) != 2 {
		return
	}

	username := parseData[0]
	lang := parseData[1]
	chatId := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID

	item, ok := b.cache.Get(username)
	if !ok {
		msg := tgbotapi.NewEditMessageText(chatId, messageID, "[информация устарела]\n\n"+update.CallbackQuery.Message.Text)
		b.bot.Request(msg)
		return
	}

	text, ok := item.data[lang]
	if !ok {
		msg := tgbotapi.NewEditMessageText(chatId, messageID, "[информация устарела]\n\n"+update.CallbackQuery.Message.Text)
		b.bot.Request(msg)
		return
	}

	msg := tgbotapi.NewEditMessageText(chatId, messageID, text)
	msg.ReplyMarkup = update.CallbackQuery.Message.ReplyMarkup
	b.bot.Request(msg)
}
