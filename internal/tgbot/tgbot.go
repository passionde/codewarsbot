package tgbot

import (
	"errors"
	"github.com/Yarik-xxx/CodeWarsRestApi/internal/app/statcollection"
	"github.com/Yarik-xxx/CodeWarsRestApi/internal/app/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"time"
)

type Bot struct {
	config *Config
	logger *logrus.Logger
	store  *statcollection.StatCollection
	bot    *tgbotapi.BotAPI
	cache  Cache
}

func New(config *Config) *Bot {
	return &Bot{
		config: config,
		logger: logrus.New(),
	}
}

func (b *Bot) Start() error {
	if err := b.configureLogger(); err != nil {
		return err
	}

	if err := b.configureStore(); err != nil {
		return err
	}

	if err := b.configureBot(); err != nil {
		return err
	}

	b.configureCache()

	updates := b.bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	for update := range updates {
		if update.CallbackQuery != nil {
			go func() {
				b.handlerSetPage(&update)
			}()
		} else if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				b.handlerStart(&update)
			default:
				go func() {
					b.handlerStatistics(&update)
				}()
			}
		}
	}
	return errors.New("stop")
}

func (b *Bot) configureStore() error {
	st := store.New(b.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	b.store = statcollection.New(st)
	return nil
}

func (b *Bot) configureLogger() error {
	level, err := logrus.ParseLevel(b.config.LogLevel)
	if err != nil {
		return err
	}

	b.logger.SetLevel(level)
	return nil
}

func (b *Bot) configureBot() error {
	bot, err := tgbotapi.NewBotAPI(b.config.Token)
	if err != nil {
		return err
	}

	b.bot = bot
	return nil
}

func (b *Bot) configureCache() {
	b.cache = NewCache()

	go func() {
		for {
			if b.cache.memory == nil {
				return
			}

			if keys := b.cache.expiredKeys(); len(keys) != 0 {
				for _, key := range keys {
					b.cache.Remove(key)
				}

			}
			time.Sleep(time.Minute)
		}
	}()
}
