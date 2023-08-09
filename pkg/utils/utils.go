package utils

import (
	"fmt"
	"itdobro/pkg/logger"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var userStatuses = make(map[int]string)

func GetBotToken() string {
	return os.Getenv("BOT_TOKEN")
}

func CreateBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.GetLogger().Panic(err)
	}

	return bot, nil
}

// SetUserStatus устанавливает статус пользователя
func SetUserStatus(userID int, status string) {
	userStatuses[userID] = status
}

// GetUserStatus получает статус пользователя
func GetUserStatus(userID int) string {
	status, exists := userStatuses[userID]
	if !exists {
		return ""
	}
	return status
}

// SendToGroup отправляет сообщение в группу
func SendToGroup(bot *tgbotapi.BotAPI, groupChatID int64, message string) {
	msg := tgbotapi.NewMessage(groupChatID, message)
	bot.Send(msg)
}

func GetUserLink(user *tgbotapi.User) string {
	return fmt.Sprintf("@%s", user.UserName)
}

func GetSocialMedia(data map[string]string) string {
	var output string

	for key, value := range data {
		output += fmt.Sprintf("[%s](%s) | ", key, value)
	}

	return "| " + output
}
