package tgbot

import (
	"fmt"
	"github.com/Yarik-xxx/CodeWarsRestApi/internal/app/statcollection"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
)

func generateMessage(u *statcollection.UserInfo) (map[string]string, []string) {
	result := make(map[string]string)
	sortKeys := []string{"all"}

	// Базовая информация
	msgBase := fmt.Sprintf("Имя: %s ", u.User.Username)

	if u.User.Name != "" {
		msgBase += fmt.Sprintf("(%s)\n", u.User.Name)
	} else {
		msgBase += "\n"
	}

	msgBase += fmt.Sprintf(
		"Баллов: %d\n"+
			"Место в рейтинге: %d\n\n"+
			"Выполнено: \nУникальных: %d\nВсего: %d\n",
		u.User.Honor,
		u.User.LeaderboardPosition,
		u.User.CodeChallenges.TotalCompletedUnique,
		u.User.CodeChallenges.TotalCompletedAll,
	)

	if u.User.CodeChallenges.TotalAuthored > 0 {
		msgBase += fmt.Sprintf("Составлено: %d\n\n", u.User.CodeChallenges.TotalAuthored)
	}

	if len(u.StatisticsKyu) == 0 {
		return map[string]string{"all": msgBase}, sortKeys
	}

	// Сводная информация об выполненных катах по рангу

	// Всего
	msgBase += "\nПо рангу:\n"

	ranks := make([]string, 0, len(u.StatisticsKyu["overall"].ByRank))
	for rank, _ := range u.StatisticsKyu["overall"].ByRank {
		ranks = append(ranks, rank)
	}
	sort.Strings(ranks)

	for _, rank := range ranks {
		nameRank := rank
		if rank == "" {
			nameRank = "Unknown"
		}
		msgBase += fmt.Sprintf("▪ %s: %d\n", nameRank, u.StatisticsKyu["overall"].ByRank[rank])
	}

	msgBase += "\nПо тегам:\n"
	// По тегам
	tags := make([]string, 0, len(u.StatisticsKyu["overall"].ByTags))
	for tag, _ := range u.StatisticsKyu["overall"].ByTags {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	for _, tag := range tags {
		msgBase += fmt.Sprintf("▪ %s: %d\n", tag, u.StatisticsKyu["overall"].ByTags[tag])
	}

	result["all"] = msgBase

	// <----------------------- По ЯП ----------------------->
	languages := make([]string, 0, len(u.StatisticsKyu))
	for lang, _ := range u.StatisticsKyu {
		if lang == "overall" {
			continue
		}
		languages = append(languages, lang)
	}
	sort.Strings(languages)
	sortKeys = append(sortKeys, languages...)

	for _, lang := range languages {
		msgCompleted := fmt.Sprintf("\n◉ %s\nВыполнено: %d\nРанг: %s\nБаллов: %d\n\nПо рангу:\n",
			lang,
			u.StatisticsKyu[lang].TotalCompleted,
			u.StatisticsKyu[lang].Rank,
			u.StatisticsKyu[lang].Score)

		// По рангу
		ranks := make([]string, 0, len(u.StatisticsKyu[lang].ByRank))
		for rank, _ := range u.StatisticsKyu[lang].ByRank {
			ranks = append(ranks, rank)
		}
		sort.Strings(ranks)

		for _, rank := range ranks {
			nameRank := rank
			if rank == "" {
				nameRank = "Unknown"
			}
			msgCompleted += fmt.Sprintf("▪ %s: %d\n", nameRank, u.StatisticsKyu[lang].ByRank[rank])
		}

		msgCompleted += "\nПо тегам:\n"
		// По тегам
		tags := make([]string, 0, len(u.StatisticsKyu[lang].ByTags))
		for tag, _ := range u.StatisticsKyu[lang].ByTags {
			tags = append(tags, tag)
		}
		sort.Strings(tags)

		for _, tag := range tags {
			msgCompleted += fmt.Sprintf("▪ %s: %d\n", tag, u.StatisticsKyu[lang].ByTags[tag])
		}

		result[lang] = msgCompleted
	}

	return result, sortKeys
}

func (b *Bot) generatePaginate(username string, keys []string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup()

	row := make([]tgbotapi.InlineKeyboardButton, 0, 3)
	for i := 0; i < len(keys); i++ {
		if i%3 == 0 && i != 0 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
			row = []tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData(
					keys[i], fmt.Sprintf("%s:%s", username, keys[i])),
			}
			continue
		}

		row = append(row, tgbotapi.NewInlineKeyboardButtonData(
			keys[i], fmt.Sprintf("%s:%s", username, keys[i])),
		)
	}

	if len(row) != 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	return keyboard
}
