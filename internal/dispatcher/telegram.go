package dispatcher

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"linnoxlewis/tg-bot-eng-dictionary/internal/config"
	"linnoxlewis/tg-bot-eng-dictionary/internal/helpers"
	"linnoxlewis/tg-bot-eng-dictionary/internal/manager"
	"linnoxlewis/tg-bot-eng-dictionary/internal/models"
	"strings"
)

const startCommand = "/start"
const translateCommand = "/t"
const generateTranslateWordsCommand = "/w"
const startLearnWord = "/learn"
const stopLearnWord = "/stop-learn-words"

var (
	errStartBot            = "start bot error"
	successStartMessage    = "Welcome to translate bot"
	successSetLearnModeBot = "Start learn mode"
	errSetLearnModeBot     = "start learn mode error"
)

type TgBotDispatcher struct {
	bot     *tgbotapi.BotAPI
	config  *config.Config
	logger  *logrus.Logger
	manager *manager.Manager
}

func NewTgBotDispatcher(config *config.Config,
	logger *logrus.Logger,
	manager *manager.Manager) *TgBotDispatcher {
	bot, err := tgbotapi.NewBotAPI(config.GetTgBotApiKey())
	if err != nil {
		logger.Panic(err)
	}
	bot.Debug = false
	logger.Printf("bot %s is working", bot.Self.UserName)

	return &TgBotDispatcher{
		bot:     bot,
		config:  config,
		logger:  logger,
		manager: manager,
	}
}

func (t *TgBotDispatcher) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5
	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		tgMsg := update.Message.Text
		msg := strings.Split(tgMsg, " ")
		command := strings.TrimSpace(msg[0])
		argument := getArgument(msg)
		userName := update.Message.From.String()
		chatId := update.Message.Chat.ID

		switch command {
		case startCommand:
			user := models.NewUser(chatId, userName)
			err := t.manager.CreateUser(ctx, user)
			if err != nil {
				t.send(chatId, errStartBot)
				break
			}
			t.send(chatId, successStartMessage)
			break

		case translateCommand:
			lang := helpers.DetectLang(argument)
			mean, err := t.manager.TranslateWordWithMeaning(argument)
			if err != nil {
				t.send(chatId, err.Error())
				break
			}
			t.send(chatId, mean.ToTranslateMessage(lang))
			break

		case startLearnWord:
			err := t.manager.UpdateLearnStatus(ctx, chatId)
			if err != nil {
				t.send(chatId, errSetLearnModeBot)
				break
			}
			t.send(chatId, successSetLearnModeBot)
			break

		case generateTranslateWordsCommand:
			words, err := t.manager.GenerateTraslateWords(t.config.GetMaxTgWord())
			if err != nil {
				t.send(chatId, err.Error())
				break
			}
			message := "Your words to learn : \n"
			for _, value := range words {
				valueString := value.GeneratingWordsToMessage()
				message = message + fmt.Sprintf("%s \n", valueString)
			}
			t.send(chatId, message)
			break

		default:
			lang := helpers.DetectLang(command)
			translate, err := t.manager.TranslateWord(command)
			if err != nil {
				t.send(chatId, err.Error())
				break
			}
			t.send(chatId, translate.ToTranslateMessage(lang))
			break
		}
	}
}

func getArgument(msg []string) string {
	argument := ""
	if len(msg) > 1 {
		argument = strings.TrimSpace(msg[1])
	}

	return argument
}

func (t *TgBotDispatcher) send(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = "html"
	_, err := t.bot.Send(msg)
	if err != nil {
		t.logger.Println(err)
	}
}
